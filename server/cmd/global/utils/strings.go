package utils

import (
	"net/url"
	"strings"
	"unicode"
)

// Check if input is a string and then check if input is an empty string
func IsEmptyString(str any) bool {
	realStr, isStr := str.(string)
	return !isStr || len(strings.TrimSpace(realStr)) <= 0
}

// Returns whether the input string contains dangerous characters or not
//
// Under the hood:
//
// return url.QueryEscape(str) == str
func IsSafeString(str string) bool {
	return url.QueryEscape(str) == str
}

// TrimSpace+ToLower+QueryEscape
func SafeString(str string) string {
	return url.QueryEscape(strings.ToLower(strings.TrimSpace(str)))
}

// Removes all "\n" from a string and TrimSpaces
func TrimString(str string) string {
	return strings.TrimSpace(strings.ReplaceAll(str, "\n", ""))
}

// Remove all non alphabetic (except "_") characters from string
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

// Remove all non alphanumerics (except "_") characters from string
func RemoveSpecialChars(str string) (out string) {
	for _, r := range str {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || unicode.IsNumber(r) || r == '_' {
			out += string(r)
		}
	}
	return
}
