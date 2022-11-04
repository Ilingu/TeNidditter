package db

import "time"

type AccountModel struct {
	AccountId uint      `json:"account_id"`
	Username  string    `json:"username"` // to index
	Password  []byte    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}
type UsersInfo struct {
	AccountModel         // extends account model
	TedditFollows string `json:"teddit_follows"`
	NitterFollows string `json:"nitter_follows"`
}

type SubtedditModel struct {
	SubID   uint   `json:"subteddit_id"`
	Subname string `json:"subname"`
}
type NittosModel struct {
	NittosID uint   `json:"twittos_id"`
	Username string `json:"username"`
}
