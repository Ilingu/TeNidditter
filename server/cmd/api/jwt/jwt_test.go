package jwt_test

import (
	"crypto/rand"
	"math/big"
	"os"
	"regexp"
	"teniditter-server/cmd/api/jwt"
	"teniditter-server/cmd/db"
	"teniditter-server/cmd/global/utils"
	utils_env "teniditter-server/cmd/global/utils/env"
	"testing"
	"time"
)

func init() {
	os.Setenv("TEST", "1")
	utils_env.LoadEnv()
}

var jwtRegex = regexp.MustCompile(`^[A-Za-z0-9-_=]+\.[A-Za-z0-9-_=]+\.?[A-Za-z0-9-_.+/=]*$`)

func TestJWT(t *testing.T) {
	username, err := utils.GenerateRandomChars(12)
	if err != nil {
		t.Fatal(err)
	}

	bigInt, err := rand.Int(rand.Reader, big.NewInt(1000))
	if err != nil {
		t.Fatal(err)
	}

	accountId, err := utils.BigIntToInt(bigInt, 16)
	if err != nil {
		t.Fatal(err)
	}

	mockAccount := db.AccountModel{AccountId: uint(accountId), Username: username}
	nowExp := time.Now().Add(time.Hour * 72).Unix()

	token, err := jwt.GenerateToken(&mockAccount)
	if err != nil {
		t.Fatal("GenerateToken() failed", err)
	}

	if !jwtRegex.MatchString(token) {
		t.Fatal("invalid jwt token")
	}

	jwtToken, err := jwt.ParseToken(token)
	if err != nil {
		t.Fatal("ParseToken() failed", err)
	}

	decodedToken, err := jwt.DecodeToken(jwtToken)
	if err != nil {
		t.Fatal("DecodeToken() failed", err)
	}

	if decodedToken.ID != mockAccount.AccountId || decodedToken.Username != mockAccount.Username || decodedToken.ExpiresAt != nowExp {
		t.Fatal("DecodeToken() decoded wrongly", decodedToken, mockAccount)
	}
}
