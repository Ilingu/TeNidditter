package utils_enc

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"teniditter-server/cmd/global/utils"
)

var iv []byte = []byte{46, 228, 83, 3, 210, 32, 229, 147, 187, 208, 189, 57, 152, 31, 7, 237}

func GenerateHashFromArgs(args ...any) string {
	concatenatedArgs := fmt.Sprint(args...)
	return Hash(concatenatedArgs)
}

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
