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

// Create a account, it returns a INCOMPLETE AccountModel based on the inputs: it only contains username, hashedPassword, and recoveryCodes. If you want the full Account created you will have to call `GetAccount`
func CreateAccount(username string, password string) (*AccountModel, error) {
	username = utils.FormatUsername(username)
	if utils.IsEmptyString(username) || len(username) < 3 || len(username) > 15 {
		return nil, errors.New("invalid username")
	}
	if !utils.IsStrongPassword(password) {
		return nil, errors.New("password too weak")
	}

	db := ps.DBManager.Connect()
	if db == nil {
		return nil, ps.ErrDbNotFound
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return nil, ErrRegister
	}

	recoveryCodes, err := generateRecoveryCodes()
	if err != nil {
		return nil, ErrRegister
	}

	hashedCodes, err := encryptRecoveryCodes(*recoveryCodes)
	if err != nil {
		return nil, ErrRegister
	}

	_, err = db.Exec("INSERT INTO Account (username, password, recovery_codes) VALUES (?, ?, ?);", username, hashedPassword, hashedCodes)
	if err != nil {
		return nil, ErrRegister
	}

	return &AccountModel{Username: username, Password: hashedPassword, RecoveryCodes: hashedCodes}, nil
}

func DeleteAccount(u *AccountModel) bool {
	db := ps.DBManager.Connect()
	if db == nil {
		return false
	}

	_, errTeship := db.Exec("DELETE FROM Teship WHERE follower_id=?;", u.AccountId)
	_, errTwiship := db.Exec("DELETE FROM Twiship WHERE follower_id=?;", u.AccountId)
	_, errLists := db.Exec("DELETE FROM NitterLists WHERE account_id=?;", u.AccountId)
	for _, err := range []error{errTeship, errTwiship, errLists} {
		if err != nil {
			return false // if one of the above query failed don't delete the account so that the user can always login
		}
	}

	_, errAccount := db.Exec("DELETE FROM Account WHERE account_id=?;", u.AccountId)
	return errAccount == nil
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
		if err := rows.Scan(&user.AccountId, &user.Username, &user.Password, &user.CreatedAt, &user.RecoveryCodes); err != nil {
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

	err := db.QueryRow("SELECT * FROM Account WHERE account_id=?", ID).Scan(&user.AccountId, &user.Username, &user.Password, &user.CreatedAt, &user.RecoveryCodes)
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

	err := db.QueryRow("SELECT * FROM Account WHERE username=?", username).Scan(&user.AccountId, &user.Username, &user.Password, &user.CreatedAt, &user.RecoveryCodes)
	if err != nil || user.AccountId == 0 || user.Username != username {
		return nil, errors.New("cannot fetch user")
	}

	return &user, nil
}
