package utils

import (
	"net/url"
	"strings"
	"unicode"

	"github.com/nbutton23/zxcvbn-go"
)

func IsEmptyString(str any) bool {
	realStr, isStr := str.(string)
	return !isStr || len(strings.TrimSpace(realStr)) <= 0
}

func FormatUsername(str string) string {
	formattedStr := ""
	for _, r := range str {
		if unicode.IsLetter(r) || r == '_' {
			formattedStr += string(r)
		}
	}

	return url.QueryEscape(strings.ToLower(strings.TrimSpace(formattedStr)))
}

func IsStrongPassword(password string) bool {
	if IsEmptyString(password) {
		return false
	}

	result := zxcvbn.PasswordStrength(password, []string{})
	return result.Score >= 3
}
