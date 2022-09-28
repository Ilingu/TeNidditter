package auth_routes

import (
	"net/http"
	"teniditter-server/cmd/api/routes"
	"teniditter-server/cmd/db"
	"teniditter-server/cmd/global/console"

	"github.com/labstack/echo/v4"
)

type RegisterPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func AuthHandler(g *echo.Group) {
	console.Log("AuthHandler Registered ✅", console.Success)

	g.POST("/register", func(c echo.Context) error {
		res := routes.EchoWrapper{Context: c}

		info := new(RegisterPayload)
		if err := c.Bind(info); err != nil {
			return res.HandleResp(http.StatusBadRequest)
		}

		account, err := db.CreateAccount(info.Username, info.Password)
		if err != nil {
			return res.HandleResp(http.StatusInternalServerError)
		}

		return res.HandleResp(http.StatusCreated, account)
	})

	g.POST("/login", func(c echo.Context) error {
		return nil
	})
}
