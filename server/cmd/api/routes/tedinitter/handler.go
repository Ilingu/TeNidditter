package tedinitter_routes

import (
	"net/http"
	"net/url"
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

		return res.HandleResp(http.StatusOK, token)
	})

	t.POST("/teddit/sub/:subname", func(c echo.Context) error {
		return SubUnsubTeddit(c, "sub")
	})

	t.DELETE("/teddit/unsub/:subname", func(c echo.Context) error {
		return SubUnsubTeddit(c, "unsub")
	})

	console.Log("TedinitterUserHandler Registered ✅", console.Info)
}

func SubUnsubTeddit(c echo.Context, method string) error {
	res := routes.EchoWrapper{Context: c}

	subname, err := url.QueryUnescape(c.Param("subname"))
	if err != nil || utils.IsEmptyString(subname) {
		return res.HandleResp(http.StatusBadRequest, "invalid subname")
	}

	token, err := jwt.DecodeToken(jwt.RetrieveToken(&c))
	if err != nil {
		return res.HandleResp(http.StatusUnauthorized, err.Error())
	}

	user := db.AccountModel{AccountId: token.ID, Username: token.Username}
	subteddit, err := db.GetSubteddit(subname)
	if err != nil {
		return res.HandleResp(http.StatusBadRequest, "failed to query this subteddit")
	}

	switch method {
	case "sub":
		if ok := user.SubToSubteddit(subteddit); !ok {
			return res.HandleResp(http.StatusInternalServerError, "couldn't subscribe you to this subteddit")
		}
	case "unsub":
		if ok := user.UnsubFromSubteddit(subteddit); !ok {
			return res.HandleResp(http.StatusInternalServerError, "couldn't unsubscibe you to this subteddit")
		}
	default:
		return res.HandleResp(http.StatusForbidden)
	}
	return res.HandleResp(http.StatusOK)
}
