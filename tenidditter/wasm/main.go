package main

import (
	"fmt"
	"syscall/js"
	"wasm/encryption"
)

func main() {
	fmt.Println("Webassembly Connected!")

	// blk, err := GenRandomBytes(16)
	// if err != nil {
	// 	panic(err)
	// }
	// iv = blk

	c := make(chan struct{}, 0)
	js.Global().Set("EncryptAES", encryption.EncryptDatas())
	js.Global().Set("DecryptAES", encryption.DecryptDatas())
	js.Global().Set("Hash", encryption.Hash())
	<-c
}
