package db

import (
	"errors"
	"log"
	"strings"
	"teniditter-server/cmd/global/utils"

	"golang.org/x/crypto/bcrypt"
)

var ErrRegister = errors.New("failed to register")

func CreateAccount(username string, password string) (*AccountModel, error) {
	db := DBManager.Connect()
	if db == nil {
		return nil, ErrDbNotFound
	}

	username = utils.FormatString(username)
	password = strings.TrimSpace(password)
	if utils.IsEmptyString(username) || utils.IsEmptyString(password) {
		return nil, ErrRegister
	}
	if len(username) < 3 || len(username) > 15 || len(password) < 15 || len(password) >= 128 {
		return nil, ErrRegister
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return nil, ErrRegister
	}

	_, err = db.Exec("INSERT INTO Account (username, password) VALUES (?, ?);", username, hashedPassword)
	if err != nil {
		return nil, ErrRegister
	}

	return nil, nil
}

func (u *AccountModel) DeleteAccount() error {
	db := DBManager.Connect()
	if db == nil {
		return ErrDbNotFound
	}

	_, err := db.Exec("DELETE FROM Account WHERE account_id=?;", u.AccountId)
	if err != nil {
		return errors.New("cannot delete account")
	}
	return nil
}

func (u *AccountModel) SignOut() bool {
	return false
}

func GetUserByUsername(username string) (*AccountModel, error) {
	db := DBManager.Connect()
	if db == nil {
		return nil, ErrDbNotFound
	}

	username = utils.FormatString(username)
	if utils.IsEmptyString(username) {
		return nil, errors.New("cannot get user")
	}

	var user AccountModel

	err := db.QueryRow("SELECT * FROM Account WHERE username LIKE ?", username).Scan(&user.AccountId, &user.Username, &user.Password, &user.CreatedAt)
	if err != nil || user.AccountId == 0 || user.Username != username {
		log.Println(err)
		return nil, errors.New("cannot fetch user")
	}

	return &user, nil
}

func (u *AccountModel) PasswordMatch(passwordInput string) bool {
	err := bcrypt.CompareHashAndPassword(u.Password, []byte(passwordInput))
	return err == nil
}
