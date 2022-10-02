package main

import (
	"crypto/aes"
	"encoding/hex"
	"fmt"
)

func main() {
	fmt.Println("Webassembly Connected!")

}

// NOTE: do not remove the "//export <name>" comments, they are here to export the funcs to wasm.exports in js

//export EncryptAES
func EncryptAES(key string, textToEnc string) string {
	c, err := aes.NewCipher([]byte(key))
	if err != nil {
		return ""
	}

	out := make([]byte, len(textToEnc))
	c.Encrypt(out, []byte(textToEnc))

	return hex.EncodeToString(out)
}

//export DecryptAES
func DecryptAES(key string, textToDec string) string {
	ciphertext, _ := hex.DecodeString(textToDec)

	c, err := aes.NewCipher([]byte(key))
	if err != nil {
		return ""
	}

	pt := make([]byte, len(ciphertext))
	c.Decrypt(pt, ciphertext)

	s := string(pt[:])
	return s
}
