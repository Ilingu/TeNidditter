package utils

import (
	crRand "crypto/rand"
	"math/big"
	"strings"
)

// It'll safely generate a true random string of "length" characters (alphanumerics, no special chars)
//
// Can be used to generate id, uuid, recovery_codes, OPT codes...
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
