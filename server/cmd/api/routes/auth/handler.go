package auth_routes

import (
	"teniditter-server/cmd/global/console"

	"github.com/labstack/echo/v4"
)

func AuthHandler(*echo.Group) {
	console.Log("AuthHandler Registered ✅", console.Success)
}
