package utils

import (
	"net/url"
	"strings"
	"unicode"
)

func IsEmptyString(str any) bool {
	realStr, isStr := str.(string)
	return !isStr || len(strings.TrimSpace(realStr)) <= 0
}

func IsSafeString(str string) bool {
	return url.QueryEscape(str) == str
}

func SafeString(str string) string {
	return url.QueryEscape(strings.ToLower(strings.TrimSpace(str)))
}

func TrimString(str string) string {
	return strings.TrimSpace(strings.ReplaceAll(str, "\n", ""))
}

func FormatString(str string) (formattedStr string) {
	for _, r := range str {
		if unicode.IsLetter(r) || r == '_' {
			formattedStr += string(r)
		}
	}
	return
}

// Remove all non alphabetic (except "_") characters from string and apply TrimSpace+ToLower+QueryEscape
func FormatUsername(str string) string {
	return SafeString(FormatString(str))
}

func RemoveSpecialChars(str string) (out string) {
	for _, r := range str {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || unicode.IsNumber(r) || r == '_' {
			out += string(r)
		}
	}
	return
}
