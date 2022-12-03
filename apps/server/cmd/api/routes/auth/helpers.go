package auth_routes

import (
	"errors"
	"teniditter-server/cmd/api/jwt"
	"teniditter-server/cmd/global/utils"

	"github.com/labstack/echo/v4"
)

func GetTokenFromQuery(c echo.Context) (*jwt.DecodedToken, error) {
	// get token from query
	tokenParams := c.QueryParam("token")
	if utils.IsEmptyString(tokenParams) {
		return nil, errors.New("missing token argument")
	}

	// parse token with the server key
	parsedToken, err := jwt.ParseToken(tokenParams)
	if err != nil {
		return nil, err
	}

	// check and decode the token to its datas
	token, err := jwt.DecodeToken(parsedToken)
	if err != nil {
		return nil, err
	}

	return &token, nil
}
