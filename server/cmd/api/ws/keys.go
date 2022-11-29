package ws

import (
	utils_enc "teniditter-server/cmd/global/utils/encryption"
)

func GenerateUserKey(ID uint, Username string) string {
	return utils_enc.GenerateHashFromArgs(ID, Username)
}
