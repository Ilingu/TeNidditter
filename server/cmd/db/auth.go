package db

func CreateAccount() (*accountModel, error) {
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
	return false
}
