package db

import (
	"errors"
	"teniditter-server/cmd/global/utils"
	ps "teniditter-server/cmd/planetscale"

	"golang.org/x/crypto/bcrypt"
)

var ErrRegister = errors.New("failed to register")

type UserInfoAcceptedArg interface {
	uint | string
}

func CreateAccount(username string, password string) (*AccountModel, error) {
	db := ps.DBManager.Connect()
	if db == nil {
		return nil, ps.ErrDbNotFound
	}

	username = utils.FormatUsername(username)
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
	db := ps.DBManager.Connect()
	if db == nil {
		return ps.ErrDbNotFound
	}

	_, err := db.Exec("DELETE FROM Account WHERE account_id=?;", u.AccountId)
	if err != nil {
		return errors.New("cannot delete account")
	}
	return nil
}

func GetAllAccounts(onlySubbedOne bool) ([]AccountModel, error) {
	db := ps.DBManager.Connect()
	if db == nil {
		return nil, ps.ErrDbNotFound
	}

	sqlQuery := "SELECT * FROM Account"
	if onlySubbedOne {
		sqlQuery = "SELECT account_id, username, password, created_at FROM Teship INNER JOIN Account ON Teship.follower_id=Account.account_id GROUP BY account_id;"
	}

	rows, err := db.Query(sqlQuery)
	if err != nil {
		return nil, errors.New("error when fetching Accounts")
	}
	defer rows.Close()

	var users []AccountModel
	for rows.Next() {
		var user AccountModel
		if err := rows.Scan(&user.AccountId, &user.Username, &user.Password, &user.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
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
	db := ps.DBManager.Connect()
	if db == nil {
		return nil, ps.ErrDbNotFound
	}

	var user AccountModel

	err := db.QueryRow("SELECT * FROM Account WHERE account_id=?", ID).Scan(&user.AccountId, &user.Username, &user.Password, &user.CreatedAt)
	if err != nil || user.AccountId == 0 || user.AccountId != ID {
		return nil, errors.New("cannot fetch user")
	}

	return &user, nil
}
func GetAccountByUsername(username string) (*AccountModel, error) {
	db := ps.DBManager.Connect()
	if db == nil {
		return nil, ps.ErrDbNotFound
	}

	username = utils.FormatUsername(username)
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
