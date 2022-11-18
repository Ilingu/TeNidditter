package utils

import (
	crRand "crypto/rand"
	"fmt"
	"math/big"
	"math/rand"
	"net/url"
	"strings"
	"time"
	"unicode"

	"teniditter-server/cmd/services/html"

	"github.com/nbutton23/zxcvbn-go"
	htmlpkg "golang.org/x/net/html"
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
func FormatUsername(str string) string {
	return SafeString(FormatString(str))
}

func RemoveSpecialChars(str string) (out string) {
	for _, r := range str {
		if (!unicode.IsMark(r) && !unicode.IsSymbol(r) && !unicode.IsPunct(r)) || r == '_' {
			out += string(r)
		}
	}
	return
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

// Will parse and tokenize the html and goes through every element of the document to check whether there is a script tag or not. If the html is invalid it'll return "true" by default
func ContainsScript(rawHtml string) (found bool) {
	err := html.FindElements(rawHtml, "script", func(elem *htmlpkg.Node) (stop bool) {
		found = true
		return true
	})
	if err != nil {
		return true
	}
	return
}

// A faster but less precise ContainsScript function: it will only look for the keyword "script" in the html, so if there is this word in a text or an attribute ect... this will return true even though there is no actual "script" tag
func ContainsScriptFast(rawHtml string) bool {
	return strings.Contains(rawHtml, "script")
}

func GenerateRandomChars(length uint) (string, error) {
	allCharSet := strings.Split("abcdedfghijklmnopqrst"+"ABCDEFGHIJKLMNOPQRSTUVWXYZ"+"0123456789", "")

	chars := []string{}
	for charId := uint(0); charId < length; charId++ {
		indexBig, err := crRand.Int(crRand.Reader, big.NewInt(int64(len(allCharSet))))
		if err != nil {
			return "", err
		}

		index, err := BigIntToInt(indexBig, 8) // bigInt must fit into int8
		if err != nil {
			return "", err
		}
		chars = append(chars, allCharSet[index])
	}
	return strings.Join(chars, ""), nil
}
