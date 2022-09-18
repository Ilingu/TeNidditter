package tedinitter_routes

import (
	"teniditter-server/cmd/global/console"

	"github.com/labstack/echo/v4"
)

func TedinitterHandler(*echo.Group) {
	console.Log("TedinitterHandler Registered ✅", console.Success)
}
