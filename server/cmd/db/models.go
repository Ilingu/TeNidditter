package db

import "time"

type AccountModel struct {
	AccountId uint      `json:"account_id"`
	Username  string    `json:"username"` // to index
	Password  []byte    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

type AccountFollows struct {
}
