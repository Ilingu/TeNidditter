package utils

import (
	"strings"
	"unicode"
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
