package tedinitter_routes

import (
	"math"
	"net/http"
	"os"
	"strconv"
	"teniditter-server/cmd/api/jwt"
	"teniditter-server/cmd/api/routes"
	nitter_routes "teniditter-server/cmd/api/routes/nitter"
	"teniditter-server/cmd/db"
	"teniditter-server/cmd/global/console"
	"teniditter-server/cmd/global/utils"
	"teniditter-server/cmd/services/nitter"

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

	/* USER */
	t.GET("/userInfo", func(c echo.Context) error {
		res := routes.EchoWrapper{Context: c}

		token, err := jwt.DecodeToken(jwt.RetrieveToken(&c))
		if err != nil {
			return res.HandleResp(http.StatusUnauthorized, err.Error())
		}

		res.SetAuthCache()
		return res.HandleResp(http.StatusOK, token)
	})
	t.PUT("/regererate-recovery-codes", func(c echo.Context) error {
		res := routes.EchoWrapper{Context: c}

		token, err := jwt.DecodeToken(jwt.RetrieveToken(&c))
		if err != nil {
			return res.HandleResp(http.StatusUnauthorized, err.Error())
		}
		user := db.AccountModel{AccountId: token.ID, Username: token.Username}

		newRecoveryCodes, err := user.RegenerateRecoveryCode()
		if err != nil {
			return res.HandleResp(http.StatusInternalServerError, err.Error())
		}

		return res.HandleResp(http.StatusOK, newRecoveryCodes)
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

	/* LISTS */

	t.POST("/nitter/list", func(c echo.Context) error {
		res := routes.EchoWrapper{Context: c}

		token, err := jwt.DecodeToken(jwt.RetrieveToken(&c))
		if err != nil {
			return res.HandleResp(http.StatusUnauthorized, err.Error())
		}
		user := db.AccountModel{AccountId: token.ID, Username: token.Username}

		type ListPayload struct {
			Listname string `json:"listname"`
		}
		payload := new(ListPayload)
		if err := c.Bind(payload); err != nil || payload == nil {
			return res.HandleResp(http.StatusBadRequest, "invalid json payload")
		}

		listname := utils.RemoveSpecialChars(payload.Listname)
		if utils.IsEmptyString(listname) || len(listname) < 3 || len(listname) >= 30 {
			return res.HandleResp(http.StatusBadRequest, "invalid listname")
		}

		if ok := db.CreateNitterList(&user, listname); !ok {
			return res.HandleResp(http.StatusInternalServerError, "failed to create list")
		}
		return res.HandleResp(http.StatusCreated)
	})
	t.DELETE("/nitter/list/:id", func(c echo.Context) error {
		res := routes.EchoWrapper{Context: c}

		token, err := jwt.DecodeToken(jwt.RetrieveToken(&c))
		if err != nil {
			return res.HandleResp(http.StatusUnauthorized, err.Error())
		}

		listId, ok := getListId(c)
		if !ok {
			return res.HandleResp(http.StatusBadRequest, "invalid list id arg")
		}

		// authorization
		list, err := db.GetListById(listId)
		if err != nil {
			return res.HandleResp(http.StatusUnauthorized, "cannot check whether you are the true owner of this list or not")
		} else if list.AccountID != token.ID {
			return res.HandleResp(http.StatusUnauthorized, "your are not the owner of this list")
		}

		if ok := db.DeleteNitterListByID(listId); !ok {
			return res.HandleResp(http.StatusInternalServerError, "failed to delete this list")
		}

		go func() {
			user := db.AccountModel{AccountId: token.ID, Username: token.Username}
			user.ListHasChange()
		}()
		return res.HandleResp(http.StatusOK)
	})

	t.GET("/nitter/lists", func(c echo.Context) error {
		res := routes.EchoWrapper{Context: c}

		token, err := jwt.DecodeToken(jwt.RetrieveToken(&c))
		if err != nil {
			return res.HandleResp(http.StatusUnauthorized, err.Error())
		}

		user := db.AccountModel{AccountId: token.ID, Username: token.Username}
		lists, err := user.GetNitterLists()
		if err != nil {
			return res.HandleResp(http.StatusInternalServerError, "Failed to fetch your lists")
		}
		return res.HandleResp(http.StatusOK, lists)
	})

	t.GET("/nitter/list/:id", func(c echo.Context) error {
		res := routes.EchoWrapper{Context: c}

		token, err := jwt.DecodeToken(jwt.RetrieveToken(&c))
		if err != nil {
			return res.HandleResp(http.StatusUnauthorized, err.Error())
		}

		listId, ok := getListId(c)
		if !ok {
			return res.HandleResp(http.StatusBadRequest, "invalid list id arg")
		}

		// authorization
		list, err := db.GetListById(listId)
		if err != nil {
			return res.HandleResp(http.StatusUnauthorized, "cannot check whether you are the true owner of this list or not")
		} else if list.AccountID != token.ID {
			return res.HandleResp(http.StatusUnauthorized, "your are not the owner of this list")
		}

		neets, err := db.GetListContentById(listId)
		if err != nil {
			return res.HandleResp(http.StatusInternalServerError, "failed to fetch the inner content")
		} else if len(neets) <= 0 {
			return res.HandleResp(http.StatusNoContent, "nothing in this list yet")
		}

		return res.HandleResp(http.StatusOK, neets)
	})
	t.POST("/nitter/list/:id/saveNeet", func(c echo.Context) error {
		res := routes.EchoWrapper{Context: c}

		token, err := jwt.DecodeToken(jwt.RetrieveToken(&c))
		if err != nil {
			return res.HandleResp(http.StatusUnauthorized, err.Error())
		}

		listId, ok := getListId(c)
		if !ok {
			return res.HandleResp(http.StatusBadRequest, "invalid list id arg")
		}

		neetPayload := new(nitter.NeetComment)
		if err := c.Bind(neetPayload); err != nil || neetPayload == nil || len(neetPayload.Id) != 19 {
			return res.HandleResp(http.StatusBadRequest, "invalid neet payload")
		}

		// authorization
		list, err := db.GetListById(listId)
		if err != nil {
			return res.HandleResp(http.StatusUnauthorized, "cannot check whether you are the true owner of this list or not")
		} else if list.AccountID != token.ID {
			return res.HandleResp(http.StatusUnauthorized, "your are not the owner of this list")
		}

		if !db.IsNeetAlreadyExist(neetPayload.Id) {
			if ok := db.InsertNewNeet(*neetPayload); !ok {
				return res.HandleResp(http.StatusInternalServerError, "couldn't add this neet to your list")
			}
		}

		if ok := list.AddNeet(neetPayload.Id); !ok {
			return res.HandleResp(http.StatusInternalServerError, "couldn't add this neet to your list")
		}
		return res.HandleResp(http.StatusOK)
	})
	t.DELETE("/nitter/list/:id/removeNeet/:neetId", func(c echo.Context) error {
		res := routes.EchoWrapper{Context: c}

		token, err := jwt.DecodeToken(jwt.RetrieveToken(&c))
		if err != nil {
			return res.HandleResp(http.StatusUnauthorized, err.Error())
		}

		listId, ok := getListId(c)
		neetId := c.Param("neetId")
		if !ok || utils.IsEmptyString(neetId) || len(neetId) != 19 {
			return res.HandleResp(http.StatusBadRequest, "invalid args")
		}

		// authorization
		list, err := db.GetListById(listId)
		if err != nil {
			return res.HandleResp(http.StatusUnauthorized, "cannot check whether you are the true owner of this list or not")
		} else if list.AccountID != token.ID {
			return res.HandleResp(http.StatusUnauthorized, "your are not the owner of this list")
		}

		if ok := list.RemoveNeet(neetId); !ok {
			return res.HandleResp(http.StatusInternalServerError, "couldn't remove this neet from your list")
		}
		return res.HandleResp(http.StatusOK)
	})

	console.Log("TedinitterUserHandler Registered ✅", console.Info)
}

func getListId(c echo.Context) (listId uint, ok bool) {
	if listidParams := c.Param("id"); !utils.IsEmptyString(listidParams) {
		if notCheckedListId, err := strconv.ParseUint(listidParams, 10, 64); err == nil {
			if notCheckedListId > 0 && notCheckedListId <= math.MaxUint {
				listId = uint(notCheckedListId)
			}
		} else {
			return 0, false
		}
	} else {
		return 0, false
	}
	return listId, true
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
			if clientIp, err := routes.GetIP(c.Request(), true); err == nil {
				go nitter_routes.StreamExternalLinks(clientIp, *feed)
			}

			return res.HandleRespBlob(http.StatusOK, feed)
		}

		// not cached: generate on the fly
		feed, err := user.GenerateNitterFeed()
		if err != nil {
			return res.HandleResp(http.StatusInternalServerError, "failed to retreive and generate this user feed: "+err.Error())
		}

		res.SetAuthCache(1800) // 30min
		if clientIp, err := routes.GetIP(c.Request(), true); err == nil {
			go nitter_routes.StreamExternalLinks(clientIp, *feed)
		}

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
