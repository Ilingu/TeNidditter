package jwt

import (
	"errors"
	"os"
	"teniditter-server/cmd/global/utils"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type JwtCustomClaims struct {
	Name  string `json:"name"`
	Admin bool   `json:"admin"`
	jwt.StandardClaims
}

func GenerateToken(username string) (string, error) {
	username = utils.FormatUsername(username)
	if utils.IsEmptyString(username) {
		return "", errors.New("invalid username, cannot generate jwt")
	}

	// Set custom claims
	claims := &JwtCustomClaims{
		username,
		false,
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
	Username  string `json:"username"`
	Admin     bool   `json:"admin"`
	ExpiresAt int64  `json:"exp"`
}

func DecodeToken(c *echo.Context) (DecodedToken, error) {
	t := (*c).Get("user").(*jwt.Token)
	claims := t.Claims.(*JwtCustomClaims)

	username := utils.FormatUsername(claims.Name)
	if utils.IsEmptyString(username) || len(username) < 3 || len(username) > 15 {
		return DecodedToken{}, errors.New("invalid username format")
	}
	if err := claims.Valid(); err != nil {
		return DecodedToken{}, errors.New("invalid token")
	}

	return DecodedToken{Username: username, Admin: false, ExpiresAt: claims.ExpiresAt}, nil
}
