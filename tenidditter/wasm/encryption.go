package main

import (
	"crypto/aes"
	"encoding/hex"
	"fmt"
	"syscall/js"
)

func main() {
	fmt.Println("Webassembly Connected!")

	c := make(chan struct{}, 0)
	js.Global().Set("EncryptAES", EncryptDatas())
	js.Global().Set("DecryptAES", DecryptDatas())
	<-c
}

func EncryptDatas() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) != 2 {
			return "Invalid no of arguments passed"
		}

		key, textToEnc := args[0].String(), args[1].String()
		return EncryptAES(key, textToEnc)
	})
}

func DecryptDatas() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) != 2 {
			return "Invalid no of arguments passed"
		}

		key, textToEnc := args[0].String(), args[1].String()
		return DecryptAES(key, textToEnc)
	})
}

func EncryptAES(key string, textToEnc string) string {
	c, err := aes.NewCipher([]byte(key))
	if err != nil {
		return ""
	}

	out := make([]byte, len(textToEnc))
	c.Encrypt(out, []byte(textToEnc))

	return hex.EncodeToString(out)
}

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
