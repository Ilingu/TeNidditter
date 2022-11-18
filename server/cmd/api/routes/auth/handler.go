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
	g.POST("/", func(c echo.Context) error {
		res := routes.EchoWrapper{Context: c}

		userInfo := new(RegisterPayload)
		if err := c.Bind(userInfo); err != nil || userInfo == nil {
			return res.HandleResp(http.StatusBadRequest, "invalid json payload")
		}
		userInfo.Username = utils.FormatUsername(userInfo.Username)

		account, err := db.GetAccount(userInfo.Username)

		if err != nil || account == nil {
			return register(res, userInfo.Username, userInfo.Password)
		}
		return login(res, account, userInfo.Password)
	})

	g.DELETE("/", func(c echo.Context) error {
		res := routes.EchoWrapper{Context: c}

		token, err := GetTokenFromQuery(c)
		if err != nil {
			return res.HandleResp(http.StatusUnauthorized, err.Error())
		}

		user := db.AccountModel{AccountId: token.ID, Username: token.Username}
		if ok := db.DeleteAccount(&user); !ok {
			return res.HandleResp(http.StatusInternalServerError, err.Error())
		}
		return logout(c)
	})

	g.GET("/available", func(c echo.Context) error {
		res := routes.EchoWrapper{Context: c}

		username := c.QueryParam("username")
		username = utils.FormatUsername(username)

		if utils.IsEmptyString(username) {
			return res.HandleResp(http.StatusBadRequest, "invalid username")
		}

		account, err := db.GetAccount(username)
		if err != nil || account == nil {
			return res.HandleResp(200, true)
		}
		return res.HandleResp(200, false)
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

	g.DELETE("/logout", logout)

	g.PUT("/reset-password", func(c echo.Context) error {
		res := routes.EchoWrapper{Context: c}

		type ResetPasswordPayload struct {
			Username     string `json:"username"`
			NewPassword  string `json:"NewPassword"`
			RecoveryCode string `json:"RecoveryCode"`
		}

		// Format datas
		resetInfo := new(ResetPasswordPayload)
		if err := c.Bind(resetInfo); err != nil || resetInfo == nil {
			return res.HandleResp(http.StatusBadRequest, "invalid json payload")
		}
		resetInfo.Username = utils.FormatUsername(resetInfo.Username)
		resetInfo.RecoveryCode = utils.TrimString(utils.RemoveSpecialChars(resetInfo.RecoveryCode))

		// Check datas
		if utils.IsEmptyString(resetInfo.Username) || len(resetInfo.Username) < 3 || len(resetInfo.Username) > 15 {
			return res.HandleResp(http.StatusBadRequest, "username not valid")
		}
		if !utils.IsStrongPassword(resetInfo.NewPassword) {
			return res.HandleResp(http.StatusBadRequest, "password not strong enough")
		}
		if utils.IsEmptyString(resetInfo.RecoveryCode) || len(resetInfo.RecoveryCode) != 8 {
			return res.HandleResp(http.StatusBadRequest, "invalid token")
		}

		// get associated user
		user, err := db.GetAccountByUsername(resetInfo.Username)
		if err != nil {
			return res.HandleResp(http.StatusNotFound, "this user is not in out database")
		}

		// check if authorized
		if validCode := user.HasRecoveryCode(resetInfo.RecoveryCode); !validCode {
			return res.HandleResp(http.StatusForbidden, "invalid token")
		}

		// consume code first, to prevent duplications
		if err := user.UseRecoveryCode(resetInfo.RecoveryCode); err != nil {
			return res.HandleResp(http.StatusInternalServerError, "cannot use this code")
		}

		// then change the user password
		if err := user.UpdatePassword(resetInfo.NewPassword); err != nil {
			// if failed revert the code by readding it
			user.AddRecoveryCode(resetInfo.RecoveryCode) // if this also fail, then f*ck it, it's just a side project
			return res.HandleResp(http.StatusInternalServerError, "couldn't update password")
		}
		return res.HandleResp(http.StatusOK)
	})

	console.Log("AuthHandler Registered ✅", console.Info)
}

func register(res routes.EchoWrapper, username, password string) error {
	if utils.IsEmptyString(username) || len(username) < 3 || len(username) > 15 {
		return res.HandleResp(http.StatusBadRequest, "username not valid")
	}
	if !utils.IsStrongPassword(password) {
		return res.HandleResp(http.StatusBadRequest, "password not strong enough")
	}

	account, err := db.CreateAccount(username, password)
	if err != nil {
		return res.HandleResp(http.StatusInternalServerError, err.Error())
	}

	if recoveryCodes, err := utils.DecryptAES(account.RecoveryCodes); err == nil {
		res.Response().Header().Set("RecoveryCodes", recoveryCodes)
	}

	return res.HandleResp(http.StatusCreated)
}

func logout(c echo.Context) error {
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
	return res.HandleResp(http.StatusResetContent)
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
	res.Response().Header()["Access-Control-Expose-Headers"] = []string{"Set-Cookie", "TedditSubs", "NitterSubs", "NitterLists"}
	Tsubs, _ := account.GetTedditSubs()
	res.InjectJsonHeader("TedditSubs", Tsubs)

	Nsubs, _ := account.GetNitterSubs()
	res.InjectJsonHeader("NitterSubs", Nsubs)

	lists, _ := account.GetNitterLists()
	res.InjectJsonHeader("NitterLists", lists)

	// adding jwt token into httpOnly cookies in the client (for future request)
	res.SetCookie(&http.Cookie{Name: "JwtToken", Value: token, Expires: time.Now().Add(30 * 24 * time.Hour), Secure: true, HttpOnly: true, SameSite: 4 /* 4 = None */, Path: "/"})

	return res.HandleResp(http.StatusAccepted, token)
}
