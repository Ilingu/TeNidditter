package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"syscall/js"
)

func main() {
	fmt.Println("Webassembly Connected!")

	// blk, err := GenRandomBytes(16)
	// if err != nil {
	// 	panic(err)
	// }
	// iv = blk

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

var iv []byte = []byte{46, 228, 83, 3, 210, 32, 229, 147, 187, 208, 189, 57, 152, 31, 7, 237}

// func GenRandomBytes(size int) (blk []byte, err error) {
// 	blk = make([]byte, size)
// 	_, err = rand.Read(blk)
// 	return
// }

func Encode(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}
func Decode(s string) []byte {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return data
}

// Encrypt method is to encrypt or hide any classified text
func EncryptAES(text, MySecret string) string {
	block, err := aes.NewCipher([]byte(MySecret))
	if err != nil {
		return ""
	}

	plainText := []byte(text)
	cfb := cipher.NewCFBEncrypter(block, iv)
	cipherText := make([]byte, len(plainText))
	cfb.XORKeyStream(cipherText, plainText)

	return Encode(cipherText)
}

// Decrypt method is to extract back the encrypted text
func DecryptAES(text, MySecret string) string {
	block, err := aes.NewCipher([]byte(MySecret))
	if err != nil {
		return ""
	}

	cipherText := Decode(text)
	cfb := cipher.NewCFBDecrypter(block, iv)
	plainText := make([]byte, len(cipherText))
	cfb.XORKeyStream(plainText, cipherText)

	return string(plainText)
}
