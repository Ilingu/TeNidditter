package db

import "time"

type AccountModel struct {
	AccountId     uint      `json:"account_id"`
	Username      string    `json:"username"` // to index
	Password      []byte    `json:"-"`
	CreatedAt     time.Time `json:"created_at"`
	RecoveryCodes string    `json:"recovery_codes"`
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

type NitterListModel struct {
	ListID    uint   `json:"list_id"`
	AccountID uint   `json:"-"`
	ListName  string `json:"title"`
}

type DBNeets struct {
	NeetId string `json:"neet_id"`
	Neet   string `json:"neet_data"`
}
