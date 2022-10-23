package ws

import (
	"fmt"
	"teniditter-server/cmd/global/utils"
)

func GenerateUserKey(ID uint, Username string) string {
	return GenerateKeyFromArgs(ID, Username)
}

func GenerateKeyFromArgs(args ...any) string {
	concatenatedArgs := fmt.Sprint(args...)
	return utils.Hash(concatenatedArgs)
}
