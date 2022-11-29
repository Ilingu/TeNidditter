package utils

import (
	"github.com/nbutton23/zxcvbn-go"
)

func IsStrongPassword(password string) bool {
	if IsEmptyString(password) {
		return false
	}

	result := zxcvbn.PasswordStrength(password, []string{})
	return result.Score >= 3
}
