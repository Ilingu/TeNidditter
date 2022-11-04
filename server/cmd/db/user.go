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
	ps "teniditter-server/cmd/planetscale"
	"teniditter-server/cmd/redis"
	"teniditter-server/cmd/redis/rediskeys"
	"teniditter-server/cmd/services/teddit"
	"time"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/websocket"
)

func (u *AccountModel) PasswordMatch(passwordInput string) bool {
	err := bcrypt.CompareHashAndPassword(u.Password, []byte(passwordInput))
	return err == nil
}

func (user AccountModel) GetTedditSubs() ([]string, error) {
	db := ps.DBManager.Connect()
	if db == nil {
		return []string{}, ps.ErrDbNotFound
	}

	rows, err := db.Query("SELECT subname FROM Teship INNER JOIN Subteddits ON Subteddits.subteddit_id = Teship.subteddit_id WHERE follower_id=?", user.AccountId)
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
	userKey := utils.GenerateKeyFromArgs(u.AccountId, u.Username /* if not enough add: u.CreatedAt.String() */)
	redisKey := rediskeys.NewKey(rediskeys.TEDDIT_USER_FEED, userKey)

	if subPosts, err := redis.Get[[]map[string]any](redisKey); err == nil {
		return &subPosts, nil // Returned from cache
	}
	return nil, errors.New("no sub posts cached for this user")
}

func (u *AccountModel) GenerateTedditFeed() (*[]map[string]any, error) {
	Tsubs, err := u.GetTedditSubs()
	if err != nil {
		return nil, err
	} else if len(Tsubs) <= 0 {
		return nil, errors.New("cannot generate feed on a user subbed to nothing")
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

	userKey := utils.GenerateKeyFromArgs(u.AccountId, u.Username /* if not enough add: u.CreatedAt.String() */)
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
			user.HasChange() // Calls ws
		}
	}()

	switch model := entity.(type) {
	case SubtedditModel:
		return user._subToSubteddit(&model, db)
	case NittosModel:
		return user._subToNittos(&model, db)
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
			user.HasChange() // Calls ws
		}
	}()

	switch model := entity.(type) {
	case SubtedditModel:
		return user._unsubFromSubteddit(&model, db)
	case NittosModel:
		return user._unsubFromNittos(&model, db)
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

/* Websocket */
type SubsPayload struct {
	Teddit []string `json:"teddit"`
	Nitter []string `json:"nitter"`
}

func (user AccountModel) HasChange() {
	wsConns, err := ws.GetWsConn(ws.GenerateUserKey(user.AccountId, user.Username))
	if err != nil || wsConns == nil {
		return
	}

	TedditSubs, err := user.GetTedditSubs()
	if err != nil {
		return
	}

	respData := SubsPayload{Teddit: TedditSubs, Nitter: []string{}}

	stringifiedSubs, err := json.Marshal(respData)
	if err != nil {
		return
	}

	for _, wsClient := range wsConns {
		go websocket.Message.Send(wsClient.WsConn, stringifiedSubs)
	}
}
