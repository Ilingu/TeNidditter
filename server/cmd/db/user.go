/*
Yes, I know this file shouldn't be under the "db" package but Go wouldn't let me split "AccountModel" into 2 pkg
so it was just unsunstainable (import cycle errors with convertions everywhere) and after literaly 1 afternoon refactoring this entire codebase I gave up and chose to split into account.go and user.go in the "db" pkg. Enventually, Go was a bad choice
*/

package db

import (
	"database/sql"
	"encoding/json"
	"errors"
	"sync"
	"teniditter-server/cmd/api/ws"
	"teniditter-server/cmd/global/utils"
	utils_enc "teniditter-server/cmd/global/utils/encryption"
	ps "teniditter-server/cmd/planetscale"
	"teniditter-server/cmd/redis"
	"teniditter-server/cmd/redis/rediskeys"
	"teniditter-server/cmd/services/nitter"
	"teniditter-server/cmd/services/teddit"
	"time"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/websocket"
)

func (u *AccountModel) PasswordMatch(passwordInput string) bool {
	err := bcrypt.CompareHashAndPassword(u.Password, []byte(passwordInput))
	return err == nil
}
func (u AccountModel) UpdatePassword(newPassword string) error {
	db := ps.DBManager.Connect()
	if db == nil {
		return ps.ErrDbNotFound
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), 12)
	if err != nil {
		return errors.New("couldn't hash new password")
	}

	_, err = db.Exec("UPDATE Account SET password=? WHERE account_id=?;", hashedPassword, u.AccountId)
	return err
}

func (user AccountModel) GetTedditSubs() ([]string, error) {
	q := "SELECT subname FROM Teship INNER JOIN Subteddits ON Subteddits.subteddit_id = Teship.subteddit_id WHERE follower_id=?"
	return user._getSubs(q)
}

func (user AccountModel) GetNitterSubs() ([]string, error) {
	q := "SELECT username FROM Twiship INNER JOIN Twittos ON Twittos.twittos_id = Twiship.twittos_id WHERE follower_id=?"
	return user._getSubs(q)
}

func (user AccountModel) _getSubs(q string) ([]string, error) {
	db := ps.DBManager.Connect()
	if db == nil {
		return []string{}, ps.ErrDbNotFound
	}

	rows, err := db.Query(q, user.AccountId)
	if err != nil {
		return []string{}, err
	}
	defer rows.Close()

	var Subs []string
	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var sub string
		if err := rows.Scan(&sub); err != nil {
			return []string{}, err
		}
		Subs = append(Subs, sub)
	}
	if err = rows.Err(); err != nil {
		return []string{}, err
	}

	return Subs, nil
}

func (u *AccountModel) GetTedditFeed() (*[]map[string]any, error) {
	userKey := utils_enc.GenerateHashFromArgs(u.AccountId, u.Username /* if not enough add: u.CreatedAt.String() */)
	redisKey := rediskeys.NewKey(rediskeys.TEDDIT_USER_FEED, userKey)

	if subPosts, err := redis.Get[[]map[string]any](redisKey); err == nil {
		return &subPosts, nil // Returned from cache
	}
	return nil, errors.New("no sub posts cached for this user")
}
func (u *AccountModel) GetNitterFeed() (*[][]nitter.NeetComment, error) {
	userKey := utils_enc.GenerateHashFromArgs(u.AccountId, u.Username /* if not enough add: u.CreatedAt.String() */)
	redisKey := rediskeys.NewKey(rediskeys.NITTER_USER_FEED, userKey)

	if subPosts, err := redis.Get[[][]nitter.NeetComment](redisKey); err == nil {
		return &subPosts, nil // Returned from cache
	}
	return nil, errors.New("no tweets cached for this user")
}

func (u *AccountModel) GenerateNitterFeed() (*[][]nitter.NeetComment, error) {
	Nsubs, err := u.GetNitterSubs()
	if err != nil {
		return nil, err
	} else if len(Nsubs) <= 0 {
		return nil, errors.New("cannot generate feed on a user subscribed to nothing")
	}

	var allTweets [][]nitter.NeetComment

	var wg sync.WaitGroup
	var mutex sync.Mutex

	wg.Add(len(Nsubs))
	for _, username := range Nsubs {
		go func(username string) {
			defer wg.Done()
			tweets, err := nitter.NittosTweetsScrap(username, 3)
			if err != nil {
				return
			}

			mutex.Lock()
			allTweets = append(allTweets, tweets...)
			mutex.Unlock()
		}(username)
	}
	wg.Wait()

	if len(allTweets) <= 0 {
		return nil, errors.New("no posts returned")
	}
	utils.ShuffleSlice(allTweets)

	userKey := utils_enc.GenerateHashFromArgs(u.AccountId, u.Username /* if not enough add: u.CreatedAt.String() */)
	redisKey := rediskeys.NewKey(rediskeys.NITTER_USER_FEED, userKey)
	go redis.Set(redisKey, allTweets, 8*time.Hour) // cache datas

	return &allTweets, nil
}
func (u *AccountModel) GenerateTedditFeed() (*[]map[string]any, error) {
	Tsubs, err := u.GetTedditSubs()
	if err != nil {
		return nil, err
	} else if len(Tsubs) <= 0 {
		return nil, errors.New("cannot generate feed on a user subscribed to nothing")
	}

	var allSubPosts []map[string]any

	var wg sync.WaitGroup
	var mutex sync.Mutex

	wg.Add(len(Tsubs))
	for _, subname := range Tsubs {
		go func(subname string) {
			defer wg.Done()
			posts, err := teddit.GetSubredditPosts(subname)
			if err != nil {
				return
			}

			linksUnknown, exist := (*posts)["links"]
			if !exist {
				return
			}

			var links []map[string]any
			blob, err := json.Marshal(linksUnknown)
			if err != nil {
				return
			}

			err = json.Unmarshal(blob, &links)
			if err != nil {
				return
			}

			mutex.Lock()
			allSubPosts = append(allSubPosts, links...)
			mutex.Unlock()
		}(subname)
	}
	wg.Wait()

	if len(allSubPosts) <= 0 {
		return nil, errors.New("no posts returned")
	}
	utils.ShuffleSlice(allSubPosts)

	userKey := utils_enc.GenerateHashFromArgs(u.AccountId, u.Username /* if not enough add: u.CreatedAt.String() */)
	redisKey := rediskeys.NewKey(rediskeys.TEDDIT_USER_FEED, userKey)
	go redis.Set(redisKey, allSubPosts, 8*time.Hour) // cache datas

	return &allSubPosts, nil
}

/* User follows */

// Sub user to a subteddit or a nittos, if successful the update will be transmitted to user via websocket
// entity is whether "SubtedditModel" or "NittosModel"
func (user AccountModel) SubTo(entity any) (success bool) {
	db := ps.DBManager.Connect()
	if db == nil {
		return false
	}

	defer func() {
		if success {
			user.SubsHasChange() // Calls ws
		}
	}()

	switch model := entity.(type) {
	case *SubtedditModel:
		return user._subToSubteddit(model, db)
	case *NittosModel:
		return user._subToNittos(model, db)
	}
	return false
}

// Unsub user from a subteddit or a nittos, if successful the update will be transmitted to user via websocket
// entity is whether "SubtedditModel" or "NittosModel"
func (user AccountModel) UnsubFrom(entity any) (success bool) {
	db := ps.DBManager.Connect()
	if db == nil {
		return false
	}

	defer func() {
		if success {
			user.SubsHasChange() // Calls ws
		}
	}()

	switch model := entity.(type) {
	case *SubtedditModel:
		return user._unsubFromSubteddit(model, db)
	case *NittosModel:
		return user._unsubFromNittos(model, db)
	}
	return false
}

// low level function that actually sub a user to a subteddit
func (user AccountModel) _subToSubteddit(sub *SubtedditModel, db *sql.DB) bool {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM Teship WHERE follower_id=? AND subteddit_id=?;", user.AccountId, sub.SubID).Scan(&count)
	if err == nil && count != 0 {
		return false
	}

	_, err = db.Exec("INSERT INTO Teship (follower_id, subteddit_id) VALUES (?,?);", user.AccountId, sub.SubID)
	return err == nil
}

// low level function that actually unsub a user from a subteddit
func (user AccountModel) _unsubFromSubteddit(sub *SubtedditModel, db *sql.DB) bool {
	_, err := db.Exec("DELETE FROM Teship WHERE follower_id=? AND subteddit_id=?;", user.AccountId, sub.SubID)
	return err == nil
}

// low level function that actually sun a user to a nittos
func (user AccountModel) _subToNittos(sub *NittosModel, db *sql.DB) bool {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM Twiship WHERE follower_id=? AND twittos_id=?;", user.AccountId, sub.NittosID).Scan(&count)
	if err == nil && count != 0 {
		return false
	}

	_, err = db.Exec("INSERT INTO Twiship (follower_id, twittos_id) VALUES (?,?);", user.AccountId, sub.NittosID)
	return err == nil
}

// low level function that actually unsub a user from a nittos
func (user AccountModel) _unsubFromNittos(sub *NittosModel, db *sql.DB) bool {
	_, err := db.Exec("DELETE FROM Twiship WHERE follower_id=? AND twittos_id=?;", user.AccountId, sub.NittosID)
	return err == nil
}

/* Lists */
func (user AccountModel) GetNitterLists() ([]NitterListModel, error) {
	db := ps.DBManager.Connect()
	if db == nil {
		return nil, ps.ErrDbNotFound
	}

	rows, err := db.Query("SELECT * FROM NitterLists WHERE account_id=?", user.AccountId)
	if err != nil {
		return nil, errors.New("error when fetching lists")
	}
	defer rows.Close()

	var lists []NitterListModel
	for rows.Next() {
		var list NitterListModel
		if err := rows.Scan(&list.ListID, &list.AccountID, &list.ListName); err != nil {
			return nil, err
		}
		lists = append(lists, list)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return lists, nil
}

/* RecoveryCodes */

func (user AccountModel) DecryptRecoveryCodes() (*[]string, error) {
	if utils.IsEmptyString(user.RecoveryCodes) {
		return nil, errors.New("no recovery codes")
	}

	blobCodes, err := utils_enc.DecryptAES(user.RecoveryCodes)
	if err != nil {
		return nil, errors.New("couldn't decrypt codes")
	}

	var recoveryCodes []string
	if err := json.Unmarshal([]byte(blobCodes), &recoveryCodes); err != nil {
		return nil, errors.New("couldn't convert codes to slice")
	}
	if len(recoveryCodes) <= 0 {
		return nil, errors.New("no recovery codes left")
	}

	return &recoveryCodes, nil
}

func (user AccountModel) AddRecoveryCode(RecoveryCode string) error {
	db := ps.DBManager.Connect()
	if db == nil {
		return ps.ErrDbNotFound
	}

	recoveryCodes, err := user.DecryptRecoveryCodes()
	if err != nil {
		return err
	}

	(*recoveryCodes) = append((*recoveryCodes), RecoveryCode)
	hashedCodes, err := encryptRecoveryCodes(*recoveryCodes)
	if err != nil {
		return errors.New("couldn't encrypt new codes")
	}
	_, err = db.Exec("UPDATE Account SET recovery_codes=? WHERE account_id=?;", hashedCodes, user.AccountId)
	return err
}

func (user AccountModel) RegenerateRecoveryCode() (*[]string, error) {
	db := ps.DBManager.Connect()
	if db == nil {
		return nil, ps.ErrDbNotFound
	}

	newRecoveryCodes, err := generateRecoveryCodes()
	if err != nil {
		return nil, errors.New("cannot regenerate codes")
	}

	hashedCodes, err := encryptRecoveryCodes(*newRecoveryCodes)
	if err != nil {
		return nil, errors.New("cannot encrypt codes")
	}

	_, err = db.Exec("UPDATE Account SET recovery_codes=? WHERE account_id=?;", hashedCodes, user.AccountId)
	if err != nil {
		return nil, err
	}

	return newRecoveryCodes, nil
}

func (user AccountModel) UseRecoveryCode(RecoveryCode string) error {
	db := ps.DBManager.Connect()
	if db == nil {
		return ps.ErrDbNotFound
	}

	recoveryCodes, err := user.DecryptRecoveryCodes()
	if err != nil {
		return err
	}

	updatedRecoveryCodes := *recoveryCodes
	for i, code := range *recoveryCodes {
		if utils_enc.Hash(code) == utils_enc.Hash(RecoveryCode) {
			updatedRecoveryCodes = append(updatedRecoveryCodes[:i], updatedRecoveryCodes[i+1:]...)
			break
		}
	}

	hashedCodes, err := encryptRecoveryCodes(updatedRecoveryCodes)
	if err != nil {
		return errors.New("couldn't encrypt new codes")
	}

	_, err = db.Exec("UPDATE Account SET recovery_codes=? WHERE account_id=?;", hashedCodes, user.AccountId)
	return err
}

func (user AccountModel) HasRecoveryCode(RecoveryCode string) bool {
	recoveryCodes, err := user.DecryptRecoveryCodes()
	if err != nil {
		return false
	}

	isValid := false
	for _, code := range *recoveryCodes {
		if utils_enc.Hash(code) == utils_enc.Hash(RecoveryCode) {
			isValid = true
			break
		}
	}
	return isValid
}

/* Websocket */
type SubsPayload struct {
	Teddit []string `json:"teddit"`
	Nitter []string `json:"nitter"`
}

func (user AccountModel) SubsHasChange() {
	wsConns, err := ws.GetWsConn(ws.GenerateUserKey(user.AccountId, user.Username))
	if err != nil || wsConns == nil {
		return
	}

	TedditSubs, _ := user.GetTedditSubs()
	NitterSubs, _ := user.GetNitterSubs()
	respData := SubsPayload{Teddit: TedditSubs, Nitter: NitterSubs}

	stringifiedSubs, err := json.Marshal(respData)
	if err != nil {
		return
	}

	for _, wsClient := range wsConns {
		go websocket.Message.Send(wsClient.WsConn, stringifiedSubs)
	}
}

func (user AccountModel) ListHasChange() {
	wsConns, err := ws.GetWsConn(ws.GenerateUserKey(user.AccountId, user.Username))
	if err != nil || wsConns == nil {
		return
	}

	NewLists, err := user.GetNitterLists()
	if err != nil {
		return
	}

	stringifiedSubs, err := json.Marshal(NewLists)
	if err != nil {
		return
	}

	for _, wsClient := range wsConns {
		go websocket.Message.Send(wsClient.WsConn, stringifiedSubs)
	}
}
