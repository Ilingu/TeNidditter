package main

import (
	"crypto/rand"
	"errors"
	"math/big"
	"strconv"
	"strings"
	"syscall/js"
)

func RandomChars() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) != 1 {
			return "Invalid no of arguments passed"
		}

		length := args[0].Int()
		randomStr, err := GenerateRandomChars(uint(length))
		if err != nil {
			return ""
		}

		return randomStr
	})
}

func GenerateRandomChars(length uint) (string, error) {
	allCharSet := strings.Split("abcdedfghijklmnopqrst"+"ABCDEFGHIJKLMNOPQRSTUVWXYZ"+"0123456789", "")

	chars := []string{}
	for charId := uint(0); charId < length; charId++ {
		indexBig, err := rand.Int(rand.Reader, big.NewInt(int64(len(allCharSet))))
		if err != nil {
			return "", err
		}

		index, err := bigIntToInt(indexBig, 8) // bigInt must fit into int8
		if err != nil {
			return "", err
		}
		chars = append(chars, allCharSet[index])
	}
	return strings.Join(chars, ""), nil
}

func bigIntToInt(bigInt *big.Int, bitSize uint8) (int64, error) {
	if bitSize > 64 {
		return 0, errors.New("invalid bitSize")
	}

	index, err := strconv.ParseInt(bigInt.String(), 10, int(bitSize)) // bigInt must fit into int8
	if err != nil {
		return 0, err
	}

	return index, nil
}
