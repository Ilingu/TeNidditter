package main

import (
	"fmt"
	"net/http"
	"os"

	auth_routes "teniditter-server/cmd/api/routes/auth"
	tedinitter_routes "teniditter-server/cmd/api/routes/tedinitter"
	"teniditter-server/cmd/db"
	"teniditter-server/cmd/global/console"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	LoadEnv() // load env if not in prod

	{
		go db.DBManager.NewDB() // Connect to DB in bg
		defer db.DBManager.Disconnect()
	}

	// Create echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.Gzip())
	enableCors(e)

	// health endpoint
	e.GET("/ping", func(c echo.Context) error {
		return c.String(200, "pong")
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
		AllowOrigins:     []string{os.Getenv("ALLOWED_ORIGIN")},
		AllowMethods:     []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
		AllowCredentials: true,
	}))

	/* To Have better security
	Content-Security-Policy: default-src 'self';base-uri 'self';font-src 'self' https: data:;form-action 'self';frame-ancestors 'self';img-src 'self' data:;object-src 'none';script-src 'self';script-src-attr 'none';style-src 'self' https: 'unsafe-inline';upgrade-insecure-requests

	Cross-Origin-Embedder-Policy: require-corp

	Cross-Origin-Opener-Policy: same-origin

	Cross-Origin-Resource-Policy: same-origin

	Origin-Agent-Cluster: ?1

	Referrer-Policy: no-referrer

	Strict-Transport-Security: max-age=15552000; includeSubDomains

	X-Content-Type-Options: nosniff

	X-DNS-Prefetch-Control: off

	X-Download-Options: noopen

	X-Frame-Options: SAMEORIGIN

	X-Permitted-Cross-Domain-Policies: none

	X-XSS-Protection: 0
	*/

	console.Log("Cors Middleware Up and Running", console.Info)
}
