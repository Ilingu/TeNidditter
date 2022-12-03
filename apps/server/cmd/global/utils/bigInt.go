package utils

import (
	"errors"
	"math/big"
	"strconv"
)

// Convert a "big.Int" to "int64" type
//
// The bitSize argument specifies the integer type that the result must fit into. Bit sizes 0, 8, 16, 32, and 64 correspond to int, int8, int16, int32, and int64. If bitSize is below 0 or above 64, an error is returned.
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
