package utils

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/rand"
	"net/url"
	"strings"
	"time"
	"unicode"

	"github.com/nbutton23/zxcvbn-go"
	"golang.org/x/net/html"
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

func TrimString(str string) string {
	return strings.TrimSpace(strings.ReplaceAll(str, "\n", ""))
}

func IsValidURL(urlToCheck string) bool {
	_, err := url.ParseRequestURI(urlToCheck)
	return err == nil
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
func FormatToSafeString(str string) string {
	return SafeString(FormatString(str))
}

func IsStrongPassword(password string) bool {
	if IsEmptyString(password) {
		return false
	}

	result := zxcvbn.PasswordStrength(password, []string{})
	return result.Score >= 3
}

func ShuffleSlice[T any](slice []T) {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(slice), func(i, j int) {
		slice[i], slice[j] = slice[j], slice[i]
	})
}

func GenerateKeyFromArgs(args ...any) string {
	concatenatedArgs := fmt.Sprint(args...)
	return Hash(concatenatedArgs)
}

func IsUrlEncoded(str string) bool {
	dec, err := url.QueryUnescape(str)
	return err != nil && dec != str
}

func ContainsScript(rawHtml string) bool {
	reader := bytes.NewReader([]byte(rawHtml))
	doc, err := html.Parse(reader)
	if err != nil {
		return true // by precaution, I return true since we cannot know if this html contains script
	}

	return findScriptInHtml(doc)
}
func findScriptInHtml(doc *html.Node) bool {
	for child := doc.FirstChild; child != nil; child = child.NextSibling {
		if child.Type == html.ElementNode && child.Data == "script" {
			return true
		}
		if found := findScriptInHtml(child); found {
			return true
		}
	}
	return false
}

func FastCheckScript(html string) bool {
	return strings.Contains(html, "script")
}
