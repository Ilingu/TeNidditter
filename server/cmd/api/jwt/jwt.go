package jwt

import (
	"errors"
	"os"
	"teniditter-server/cmd/db"
	"teniditter-server/cmd/global/utils"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type JwtCustomClaims struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	// TedditSubs []string `json:"teddit_subs"`
	jwt.StandardClaims
}

func GenerateToken(account *db.AccountModel) (string, error) {
	username := utils.FormatToSafeString(account.Username)
	if utils.IsEmptyString(username) {
		return "", errors.New("invalid username, cannot generate jwt")
	}

	// Set custom claims
	claims := &JwtCustomClaims{
		account.AccountId,
		username,
		// tedditSubs,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	if len(os.Getenv("JWT_SECRET")) < 15 {
		return "", errors.New("invalid secret, cannot generate jwt")
	}

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}

	return t, nil
}

type DecodedToken struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	ExpiresAt int64  `json:"exp"`
}

func DecodeToken(c *echo.Context) (DecodedToken, error) {
	t := (*c).Get("user").(*jwt.Token)
	claims := t.Claims.(*JwtCustomClaims)

	username := utils.FormatToSafeString(claims.Name)
	if utils.IsEmptyString(username) || len(username) < 3 || len(username) > 15 {
		return DecodedToken{}, errors.New("invalid username format")
	}
	if err := claims.Valid(); err != nil {
		return DecodedToken{}, errors.New("invalid token")
	}

	return DecodedToken{claims.ID, username, claims.ExpiresAt}, nil
}
