package utils_enc

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"teniditter-server/cmd/global/utils"
)

// Concatenate the args into a string and then returns the hash of this string
func GenerateHashFromArgs(args ...any) string {
	concatenatedArgs := fmt.Sprint(args...)
	return Hash(concatenatedArgs)
}

// Simple sha256 hash function
func Hash(str string) string {
	ByteHash := sha256.Sum256([]byte(str))
	HashedStr := fmt.Sprintf("%x", ByteHash[:])
	return HashedStr
}

// Encrypt method is to encrypt or hide any classified text
func EncryptAES(textToEnc string) (string, error) {
	if utils.IsEmptyString(os.Getenv("ENCRYPTION_KEY")) {
		return "", errors.New("no enc key")
	}

	iv, err := getIV()
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher([]byte(os.Getenv("ENCRYPTION_KEY")))
	if err != nil {
		return "", err
	}

	plainText := []byte(textToEnc)
	cfb := cipher.NewCFBEncrypter(block, iv)
	cipherText := make([]byte, len(plainText))
	cfb.XORKeyStream(cipherText, plainText)

	return encode(cipherText), nil
}

// Decrypt method is to extract back the encrypted text
func DecryptAES(textToDec string) (string, error) {
	if utils.IsEmptyString(os.Getenv("ENCRYPTION_KEY")) {
		return "", errors.New("no enc key")
	}

	iv, err := getIV()
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher([]byte(os.Getenv("ENCRYPTION_KEY")))
	if err != nil {
		return "", err
	}

	cipherText, err := decode(textToDec)
	if err != nil {
		return "", err
	}

	cfb := cipher.NewCFBDecrypter(block, iv)
	plainText := make([]byte, len(cipherText))
	cfb.XORKeyStream(plainText, cipherText)

	return string(plainText), nil
}

// Helpers

// get IV from env and Unmarshal it into "[]byte"
func getIV() ([]byte, error) {
	if utils.IsEmptyString(os.Getenv("IV_KEY")) {
		return nil, errors.New("no iv key")
	}

	var iv []byte
	if err := json.Unmarshal([]byte(os.Getenv("IV_KEY")), &iv); err != nil {
		return nil, errors.New("no iv key")
	}

	return iv, nil
}
func encode(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}
func decode(s string) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}
	return data, nil
}
