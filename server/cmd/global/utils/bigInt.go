package utils

import (
	"errors"
	"math/big"
	"strconv"
)

func BigIntToInt(bigInt *big.Int, bitSize uint8) (int64, error) {
	if bitSize > 64 {
		return 0, errors.New("invalid bitSize")
	}

	index, err := strconv.ParseInt(bigInt.String(), 10, int(bitSize)) // bigInt must fit into int8
	if err != nil {
		return 0, err
	}

	return index, nil
}
