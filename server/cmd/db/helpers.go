package db

import (
	"encoding/json"
	"errors"
	"teniditter-server/cmd/global/utils"
)

func encryptRecoveryCodes(recoveryCodes []string) (string, error) {
	blobCodes, err := json.Marshal(recoveryCodes)
	if err != nil {
		return "", errors.New("couldn't marshal codes")
	}

	hashedCodes, err := utils.EncryptAES(string(blobCodes))
	if err != nil {
		return "", errors.New("couldn't encrypt codes")
	}

	return hashedCodes, nil
}

var (
	RECOVERY_CODES_AMOUNT = 6
	RECOVERY_CODE_LENGTH  = 10
)

func generateRecoveryCodes() (*[]string, error) {
	recoveryCodes := []string{}
	for nGen := 0; ; nGen++ {
		if nGen > 100 {
			return nil, errors.New("couldn't generate user recovery codes")
		}

		tempCodes := []string{}
		for i := 0; i < RECOVERY_CODES_AMOUNT; i++ {
			code, err := generateRecoveryCode()
			if err != nil {
				break
			}
			tempCodes = append(tempCodes, code)
		}

		if len(tempCodes) == RECOVERY_CODES_AMOUNT {
			recoveryCodes = tempCodes
			break
		}
	}
	return &recoveryCodes, nil
}

func generateRecoveryCode() (string, error) {
	code, err := utils.GenerateRandomChars(uint(RECOVERY_CODE_LENGTH))
	if err != nil {
		return "", err
	}
	return code, nil
}
