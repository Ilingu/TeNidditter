package utils

import (
	"strings"
	"unicode"

	"github.com/nbutton23/zxcvbn-go"
)

func IsEmptyString(str any) bool {
	realStr, isStr := str.(string)
	return !isStr || len(strings.TrimSpace(realStr)) <= 0
}

func FormatString(str string) string {
	formattedStr := ""
	for _, r := range str {
		if unicode.IsLetter(r) {
			formattedStr += string(r)
		}
	}

	return strings.ToLower(strings.TrimSpace(formattedStr))
}

func IsStrongPassword(password string) bool {
	if IsEmptyString(password) {
		return false
	}

	result := zxcvbn.PasswordStrength(password, []string{})
	return result.Score >= 3
}
