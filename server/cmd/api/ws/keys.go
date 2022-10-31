package ws

import (
	"teniditter-server/cmd/global/utils"
)

func GenerateUserKey(ID uint, Username string) string {
	return utils.GenerateKeyFromArgs(ID, Username)
}
