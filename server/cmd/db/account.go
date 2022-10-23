package db

import (
	"errors"
	"teniditter-server/cmd/global/utils"

	"golang.org/x/crypto/bcrypt"
)

var ErrRegister = errors.New("failed to register")

type UserInfoAcceptedArg interface {
	uint | string
}

func CreateAccount(username string, password string) (*AccountModel, error) {
	db := DBManager.Connect()
	if db == nil {
		return nil, ErrDbNotFound
	}

	username = utils.FormatToSafeString(username)
	if utils.IsEmptyString(username) || len(username) < 3 || len(username) > 15 {
		return nil, errors.New("invalid username")
	}
	if !utils.IsStrongPassword(password) {
		return nil, errors.New("password too weak")
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

func DeleteAccount(u *AccountModel) error {
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

// Get User by ID or username
func GetAccount[T UserInfoAcceptedArg](username_or_userId T) (*AccountModel, error) {
	switch realVal := any(username_or_userId).(type) {
	case uint:
		return GetAccountByID(realVal)
	case string:
		return GetAccountByUsername(realVal)
	default:
		return nil, errors.New("invalid user info (type: uint (userID) or string (username))")
	}
}
func GetAccountByID(ID uint) (*AccountModel, error) {
	db := DBManager.Connect()
	if db == nil {
		return nil, ErrDbNotFound
	}

	var user AccountModel

	err := db.QueryRow("SELECT * FROM Account WHERE account_id=?", ID).Scan(&user.AccountId, &user.Username, &user.Password, &user.CreatedAt)
	if err != nil || user.AccountId == 0 || user.AccountId != ID {
		return nil, errors.New("cannot fetch user")
	}

	return &user, nil
}
func GetAccountByUsername(username string) (*AccountModel, error) {
	db := DBManager.Connect()
	if db == nil {
		return nil, ErrDbNotFound
	}

	username = utils.FormatToSafeString(username)
	if utils.IsEmptyString(username) {
		return nil, errors.New("cannot get user")
	}

	var user AccountModel

	err := db.QueryRow("SELECT * FROM Account WHERE username=?", username).Scan(&user.AccountId, &user.Username, &user.Password, &user.CreatedAt)
	if err != nil || user.AccountId == 0 || user.Username != username {
		return nil, errors.New("cannot fetch user")
	}

	return &user, nil
}
