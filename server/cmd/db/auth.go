package db

import (
	"errors"
	"strings"
	"teniditter-server/cmd/global/utils"

	"golang.org/x/crypto/bcrypt"
)

var ErrRegister = errors.New("failed to register")

func CreateAccount(username string, password string) (*accountModel, error) {
	db := DBManager.Connect()
	if db == nil {
		return nil, ErrDbNotFound
	}

	username = strings.TrimSpace(username)
	password = strings.TrimSpace(password)
	if utils.IsEmptyString(username) || utils.IsEmptyString(password) || len(password) <= 15 {
		return nil, ErrRegister
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return nil, ErrRegister
	}

	insert, err := db.Query("INSERT INTO Account (username, password) VALUES (?, ?)", username, hashedPassword)
	if err != nil {
		return nil, ErrRegister
	}
	defer insert.Close()

	return nil, nil
}

func (u *accountModel) DeleteAccount() error {
	return nil
}

func (u *accountModel) SignOut() bool {
	return false
}

func GetUserByID() *accountModel {
	return nil
}

func (u *accountModel) PasswordMatch(passwordInput string) bool {
	err := bcrypt.CompareHashAndPassword(u.Password, []byte(passwordInput))
	return err == nil
}
