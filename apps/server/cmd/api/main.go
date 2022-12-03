package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	auth_routes "teniditter-server/cmd/api/routes/auth"
	cron_routes "teniditter-server/cmd/api/routes/cron"
	nitter_routes "teniditter-server/cmd/api/routes/nitter"
	teddit_routes "teniditter-server/cmd/api/routes/teddit"
	tedinitter_routes "teniditter-server/cmd/api/routes/tedinitter"
	"teniditter-server/cmd/global/console"
	utils_env "teniditter-server/cmd/global/utils/env"
	ps "teniditter-server/cmd/planetscale"
	"teniditter-server/cmd/redis"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	utils_env.LoadEnv() // load env if not in prod

	if os.Getenv("TEST") != "1" {
		go ps.DBManager.NewDB() // Connect to DB in bg
		defer ps.DBManager.Disconnect()

		go redis.ConnectRedis()
		defer redis.DisconnectRedis()
	}

	// Create echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))
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
	tedinitter_routes.TedinitterUserHandler(tedinitterG)

	tedditG := e.Group("/teddit")
	teddit_routes.TedditHandler(tedditG)

	nitterG := e.Group("/nitter")
	nitter_routes.NitterHandler(nitterG)

	cronG := e.Group("/cron")
	cron_routes.CronListener(cronG)
	// go cron_routes.RegisterCron()

	// Start Server
	PORT := fmt.Sprintf(":%s", os.Getenv("PORT"))

	if os.Getenv("TEST") == "1" {
		timeout := time.After(5 * time.Second)
		go e.Start(PORT)

		<-timeout
		return
	}

	if err := e.Start(PORT); err != nil {
		e.Close()
		e.Logger.Fatal("failed to start the server", err)
	}
}

func enableCors(e *echo.Echo) {
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{os.Getenv("ALLOWED_ORIGIN")},
		AllowMethods:     []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
		ExposeHeaders:    []string{"TedditSubs"},
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

	console.Log("Cors Middleware Up and Running")
}
