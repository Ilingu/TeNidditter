package tedinitter_routes

import (
	"net/http"
	"os"
	"teniditter-server/cmd/api/jwt"
	"teniditter-server/cmd/api/routes"
	"teniditter-server/cmd/global/console"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func TedinitterHandler(t *echo.Group) {
	if len(os.Getenv("JWT_SECRET")) < 15 {
		console.Log("Couldn't register Tedinitter routes: JWT_SECRET is not secured", console.Error)
		return
	}

	config := middleware.JWTConfig{
		Claims:     &jwt.JwtCustomClaims{},
		SigningKey: []byte(os.Getenv("JWT_SECRET")),
	}
	t.Use(middleware.JWTWithConfig(config)) // restricted routes

	t.GET("/userInfo", func(c echo.Context) error {
		res := routes.EchoWrapper{Context: c}

		token, err := jwt.DecodeToken(&c)
		if err != nil {
			return res.HandleResp(http.StatusUnauthorized, err.Error())
		}

		return res.HandleResp(http.StatusOK, token)
	})

	console.Log("TedinitterHandler Registered ✅", console.Success)
}
