/*
Yes, I know this file shouldn't be under the "db" package but Go wouldn't let me split "AccountModel" into 2 pkg
so it was just unsunstainable (import cycle errors with convertions everywhere) and after literaly 1 afternoon refactoring this entire codebase I gave up and chose to split into account.go and user.go in the "db" pkg. Enventually, Go was a bad choice
*/

package db

import (
	"encoding/json"
	"teniditter-server/cmd/api/ws"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/websocket"
)

func (u *AccountModel) PasswordMatch(passwordInput string) bool {
	err := bcrypt.CompareHashAndPassword(u.Password, []byte(passwordInput))
	return err == nil
}

func (user AccountModel) GetTedditSubs() ([]string, error) {
	db := DBManager.Connect()
	if db == nil {
		return []string{}, ErrDbNotFound
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

func (user AccountModel) SubToSubteddit(sub *SubtedditModel) bool {
	db := DBManager.Connect()
	if db == nil {
		return false
	}

	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM Teship WHERE follower_id=? AND subteddit_id=?;", user.AccountId, sub.SubID).Scan(&count)
	if err == nil && count != 0 {
		return false
	}

	_, err = db.Exec("INSERT INTO Teship (follower_id, subteddit_id) VALUES (?,?);", user.AccountId, sub.SubID)
	if err != nil {
		return false
	}

	user.HasChange() // Calls ws
	return true
}

func (user AccountModel) UnsubFromSubteddit(sub *SubtedditModel) bool {
	db := DBManager.Connect()
	if db == nil {
		return false
	}

	_, err := db.Exec("DELETE FROM Teship WHERE follower_id=? AND subteddit_id=?;", user.AccountId, sub.SubID)
	if err != nil {
		return false
	}

	user.HasChange() // Calls ws
	return true
}

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

	for _, wsClient := range *wsConns {
		go websocket.Message.Send(wsClient.WsConn, stringifiedSubs)
	}
}
