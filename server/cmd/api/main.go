package main

import (
	"fmt"
	"net/http"
	"os"

	auth_routes "teniditter-server/cmd/api/routes/auth"
	tedinitter_routes "teniditter-server/cmd/api/routes/tedinitter"
	"teniditter-server/cmd/global/console"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	LoadEnv() // load env if not in prod

	// Create echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.Gzip())
	enableCors(e)

	// health endpoint
	e.GET("/ping", func(c echo.Context) error {
		return c.File("../../go.sum")
	})

	// Subroutes
	authG := e.Group("/auth")
	auth_routes.AuthHandler(authG)

	tedinitterG := e.Group("/tedinitter")
	tedinitter_routes.TedinitterHandler(tedinitterG)

	// Start Server
	PORT := fmt.Sprintf(":%s", os.Getenv("PORT"))
	e.Logger.Fatal(e.Start(PORT))
}

func enableCors(e *echo.Echo) {
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{os.Getenv("ALLOWED_ORIGIN")},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))

	console.Log("Cors Middleware Up and Running", console.Info)
}
