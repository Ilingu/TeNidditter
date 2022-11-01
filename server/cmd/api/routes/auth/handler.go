package auth_routes

import (
	"net/http"
	"teniditter-server/cmd/api/jwt"
	"teniditter-server/cmd/api/routes"
	"teniditter-server/cmd/api/ws"
	"teniditter-server/cmd/db"
	"teniditter-server/cmd/global/console"
	"teniditter-server/cmd/global/utils"
	"time"

	"github.com/labstack/echo/v4"
	"golang.org/x/net/websocket"
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

		account, err := db.GetAccount(userInfo.Username)

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

		account, err := db.GetAccount(username)
		if err != nil || account == nil {
			return res.HandleResp(200, true)
		}
		return res.HandleResp(200, false)
	})

	g.DELETE("/erase", func(c echo.Context) error {
		res := routes.EchoWrapper{Context: c}
		return res.HandleResp(http.StatusNotImplemented, "Not Implemented Yet")
	})

	g.GET("/userChanged", func(c echo.Context) error {
		res := routes.EchoWrapper{Context: c}

		token, err := GetTokenFromQuery(c)
		if err != nil {
			return res.HandleResp(http.StatusUnauthorized, err.Error())
		}

		cws := ws.EchoWrapper{Context: c}
		cws.NewWsConn(ws.GenerateUserKey(token.ID, token.Username), make(chan *ws.WebsocketConn)) // reply to the caller
		return nil
	})

	g.DELETE("/logout", func(c echo.Context) error {
		res := routes.EchoWrapper{Context: c}

		// Close all ws connection related to this user
		if token, err := GetTokenFromQuery(c); err == nil {
			if wsConns, err := ws.GetWsConn(ws.GenerateUserKey(token.ID, token.Username)); err == nil {
				for _, conn := range wsConns {
					websocket.Message.Send(conn.WsConn, "LOGOUT") // logout user of all others client
					conn.CloseConn()                              // closing ws conn on server
				}
			}
		}

		res.Response().Header().Set("Clear-Site-Data", `"cache", "cookies", "storage", "executionContexts"`)
		res.Response().Header().Set("Access-Control-Expose-Headers", "Clear-Site-Data")
		return res.HandleResp(205)
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

	token, err := jwt.GenerateToken(account)
	if err != nil {
		return res.HandleResp(http.StatusInternalServerError, "Couldn't Generate JWT token")
	}

	// custom headers
	res.Response().Header()["Access-Control-Expose-Headers"] = []string{"Set-Cookie", "TedditSubs"}
	subs, _ := account.GetTedditSubs()
	res.InjectSubs(subs)

	// adding jwt token into httpOnly cookies in the client (for future request)
	res.SetCookie(&http.Cookie{Name: "JwtToken", Value: token, Expires: time.Now().Add(30 * 24 * time.Hour), Secure: true, HttpOnly: true, SameSite: 4 /* 4 = None */, Path: "/"})

	return res.HandleResp(http.StatusAccepted, token)
}
