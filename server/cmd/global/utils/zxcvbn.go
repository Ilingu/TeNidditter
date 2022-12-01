package utils

import (
	"github.com/nbutton23/zxcvbn-go"
)

// Returns whether the input password meet the app's security standars
//
// Usage:
//
// IsStrongPassword("1234") // false
//
// IsStrongPassword("os$8@2w%XJ$m8V0MAY0icN#4Yd2tka6L") // true
func IsStrongPassword(password string) bool {
	if IsEmptyString(password) {
		return false
	}

	result := zxcvbn.PasswordStrength(password, []string{})
	return result.Score >= 3
}
