package utils

import (
	"crypto/sha256"
	"fmt"
	"net/url"
	"strings"
	"unicode"

	"github.com/nbutton23/zxcvbn-go"
)

func Hash(str string) string {
	ByteHash := sha256.Sum256([]byte(str))
	HashedStr := fmt.Sprintf("%x", ByteHash[:])
	return HashedStr
}

func IsEmptyString(str any) bool {
	realStr, isStr := str.(string)
	return !isStr || len(strings.TrimSpace(realStr)) <= 0
}

func SafeString(str string) string {
	return url.QueryEscape(strings.ToLower(strings.TrimSpace(str)))
}

func IsValidURL(urlToCheck string) bool {
	_, err := url.ParseRequestURI(urlToCheck)
	return err == nil
}

// Remove all non alphabetic (except "_") characters from string and apply TrimSpace+ToLower+QueryEscape
func FormatToSafeString(str string) string {
	formattedStr := ""
	for _, r := range str {
		if unicode.IsLetter(r) || r == '_' {
			formattedStr += string(r)
		}
	}

	return SafeString(formattedStr)
}

func IsStrongPassword(password string) bool {
	if IsEmptyString(password) {
		return false
	}

	result := zxcvbn.PasswordStrength(password, []string{})
	return result.Score >= 3
}
