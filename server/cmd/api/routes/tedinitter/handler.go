package tedinitter_routes

import (
	"net/http"
	"os"
	"teniditter-server/cmd/api/jwt"
	"teniditter-server/cmd/api/routes"
	"teniditter-server/cmd/db"
	"teniditter-server/cmd/global/console"
	"teniditter-server/cmd/global/utils"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func TedinitterUserHandler(t *echo.Group) {
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

		token, err := jwt.DecodeToken(jwt.RetrieveToken(&c))
		if err != nil {
			return res.HandleResp(http.StatusUnauthorized, err.Error())
		}

		res.SetAuthCache()
		return res.HandleResp(http.StatusOK, token)
	})

	/* TEDDIT */

	t.POST("/teddit/sub/:name", func(c echo.Context) error {
		return handleSubUnsub(c, "teddit", "sub")
	})

	t.DELETE("/teddit/unsub/:name", func(c echo.Context) error {
		return handleSubUnsub(c, "teddit", "unsub")
	})

	t.GET("/teddit/feed", func(c echo.Context) error {
		return handleGetFeed(c, "teddit")
	})

	/* NITTER */
	t.POST("/nitter/sub/:name", func(c echo.Context) error {
		return handleSubUnsub(c, "nitter", "sub")
	})

	t.DELETE("/nitter/unsub/:name", func(c echo.Context) error {
		return handleSubUnsub(c, "nitter", "unsub")
	})

	t.GET("/nitter/feed", func(c echo.Context) error {
		return handleGetFeed(c, "nitter")
	})

	console.Log("TedinitterUserHandler Registered ✅", console.Info)
}

func handleGetFeed(c echo.Context, service string) error {
	res := routes.EchoWrapper{Context: c}

	token, err := jwt.DecodeToken(jwt.RetrieveToken(&c))
	if err != nil {
		return res.HandleResp(http.StatusUnauthorized, err.Error())
	}

	user := db.AccountModel{AccountId: token.ID, Username: token.Username}

	switch service {
	case "teddit":
		if feed, err := user.GetTedditFeed(); err == nil {
			console.Log("Feed Returned from cache", console.Neutral)
			res.SetAuthCache(1800) // 30min
			return res.HandleRespBlob(http.StatusOK, feed)
		}

		// not cached: generate on the fly
		feed, err := user.GenerateTedditFeed()
		if err != nil {
			return res.HandleResp(http.StatusInternalServerError, "failed to retreive and generate this user feed: "+err.Error())
		}
		res.SetAuthCache(1800) // 30min
		return res.HandleRespBlob(http.StatusOK, feed)
	case "nitter":
		if feed, err := user.GetNitterFeed(); err == nil {
			console.Log("Feed Returned from cache", console.Neutral)
			res.SetAuthCache(1800) // 30min
			return res.HandleRespBlob(http.StatusOK, feed)
		}

		// not cached: generate on the fly
		feed, err := user.GenerateNitterFeed()
		if err != nil {
			return res.HandleResp(http.StatusInternalServerError, "failed to retreive and generate this user feed: "+err.Error())
		}
		res.SetAuthCache(1800) // 30min
		return res.HandleRespBlob(http.StatusOK, feed)
	default:
		return res.HandleRespBlob(http.StatusNotAcceptable, "service not found")
	}
}

// service is whether "teddit" or "nitter" and action is whether "sub" or "unsub"
func handleSubUnsub(c echo.Context, service, action string) error {
	res := routes.EchoWrapper{Context: c}

	entityName := c.Param("name")
	if utils.IsEmptyString(entityName) {
		return res.HandleResp(http.StatusBadRequest, "invalid name")
	}

	token, err := jwt.DecodeToken(jwt.RetrieveToken(&c))
	if err != nil {
		return res.HandleResp(http.StatusUnauthorized, err.Error())
	}

	user := db.AccountModel{AccountId: token.ID, Username: token.Username}

	var entity any
	switch service {
	case "teddit":
		if subteddit, err := db.GetSubteddit(entityName); err == nil {
			entity = subteddit
		} else {
			return res.HandleResp(http.StatusBadRequest, "failed to query this subteddit")
		}
	case "nitter":
		if nittos, err := db.GetNittos(entityName); err == nil {
			entity = nittos
		} else {
			return res.HandleResp(http.StatusBadRequest, "failed to query this nittos")
		}
	default:
		return res.HandleResp(http.StatusBadRequest, "invalid service")
	}

	switch action {
	case "sub":
		if ok := user.SubTo(entity); !ok {
			return res.HandleResp(http.StatusInternalServerError, "couldn't subscribe you to this entity")
		}
	case "unsub":
		if ok := user.UnsubFrom(entity); !ok {
			return res.HandleResp(http.StatusInternalServerError, "couldn't unsubscibe you from this entity")
		}
	default:
		return res.HandleResp(http.StatusForbidden)
	}
	return res.HandleResp(http.StatusOK)
}
