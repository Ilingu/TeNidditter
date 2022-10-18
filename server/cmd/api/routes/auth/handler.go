package auth_routes

import (
	"net/http"
	"teniditter-server/cmd/api/jwt"
	"teniditter-server/cmd/api/routes"
	"teniditter-server/cmd/db"
	"teniditter-server/cmd/global/console"
	"teniditter-server/cmd/global/utils"

	"github.com/labstack/echo/v4"
)

type RegisterPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func AuthHandler(g *echo.Group) {
	console.Log("AuthHandler Registered ✅", console.Info)

	g.POST("/", func(c echo.Context) error {
		res := routes.EchoWrapper{Context: c}

		userInfo := new(RegisterPayload)
		if err := c.Bind(userInfo); err != nil {
			return res.HandleResp(http.StatusBadRequest)
		}

		account, err := db.GetUserByUsername(userInfo.Username)

		if err != nil || account == nil {
			return register(res, userInfo.Username, userInfo.Password)
		}
		return login(res, account, userInfo.Password)
	})

	g.GET("/available", func(c echo.Context) error {
		res := routes.EchoWrapper{Context: c}

		username := c.QueryParam("username")
		username = utils.FormatToSafeString(username)

		if utils.IsEmptyString(username) {
			return res.HandleResp(http.StatusBadRequest, "invalid username")
		}

		account, err := db.GetUserByUsername(username)
		if err != nil || account == nil {
			return res.HandleResp(200, true)
		}
		return res.HandleResp(200, false)
	})

	g.DELETE("/erase", func(c echo.Context) error {
		res := routes.EchoWrapper{Context: c}
		return res.HandleResp(http.StatusNotImplemented, "Not Implemented Yet")
	})
}

func register(res routes.EchoWrapper, username, password string) error {
	account, err := db.CreateAccount(username, password)
	if err != nil {
		return res.HandleResp(http.StatusInternalServerError, err.Error())
	}

	return res.HandleResp(http.StatusCreated, *account)
}

func login(res routes.EchoWrapper, account *db.AccountModel, password string) error {
	if authenticated := account.PasswordMatch(password); !authenticated {
		return res.HandleResp(http.StatusForbidden, "Wrong Credentials")
	}

	token, err := jwt.GenerateToken(account.Username)
	if err != nil {
		return res.HandleResp(http.StatusInternalServerError, "Couldn't Generate JWT token")
	}

	return res.HandleResp(http.StatusAccepted, token)
}
