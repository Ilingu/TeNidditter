package nitter_routes

import (
	"net/http"
	"strconv"
	"teniditter-server/cmd/api/routes"
	"teniditter-server/cmd/global/console"
	"teniditter-server/cmd/global/utils"
	"teniditter-server/cmd/services/nitter"

	"github.com/labstack/echo/v4"
)

func NitterHandler(n *echo.Group) {
	n.GET("/search", func(c echo.Context) error {
		res := routes.EchoWrapper{Context: c}

		QuerySearch := c.QueryParam("q")
		SearchType := c.QueryParam("type")
		if utils.IsEmptyString(QuerySearch) || utils.IsEmptyString(SearchType) {
			return res.HandleResp(http.StatusBadRequest, `invalid query "type" or "q"`)
		}

		queryLimit := 1
		if limit, err := strconv.ParseInt(c.QueryParam("limit"), 10, 8); err == nil && limit > 0 && limit <= 127 {
			queryLimit = int(limit)
		}

		switch SearchType {
		case "tweets":
			tweets, err := nitter.SearchTweetsScrap(QuerySearch, queryLimit)
			if err != nil {
				return res.HandleResp(http.StatusNotFound, err.Error())
			}

			res.SetPublicCache(15 * 60) // 15min
			return res.HandleResp(http.StatusOK, tweets)

		case "users":
			if !utils.IsSafeString(QuerySearch) {
				return res.HandleResp(http.StatusBadRequest, `invalid username`)
			}

			nittos, err := nitter.SearchNittos(QuerySearch, queryLimit)
			if err != nil {
				return res.HandleResp(http.StatusNotFound, err.Error())
			}

			res.SetPublicCache(15 * 60) // 15min
			return res.HandleResp(http.StatusOK, nittos)
		default:
			return res.HandleResp(http.StatusBadRequest, `invalid query "type", it must be whether "tweets" or "users"`)
		}
	})

	n.GET("/nittos/:name/about", func(c echo.Context) error {
		res := routes.EchoWrapper{Context: c}

		username := c.Param("name")
		if utils.IsEmptyString(username) {
			return res.HandleResp(http.StatusBadRequest, "invalid username param")
		}

		metadata, err := nitter.NittosMetadata(username)
		if err != nil {
			return res.HandleResp(http.StatusNotFound, "no metadata returned for this user")
		}

		res.SetPublicCache(2 * 24 * 60 * 60) // 2d
		return res.HandleResp(http.StatusOK, metadata)
	})
	n.GET("/nittos/:name/neets", func(c echo.Context) error {
		res := routes.EchoWrapper{Context: c}

		username := (c.Param("name"))
		if utils.IsEmptyString(username) {
			return res.HandleResp(http.StatusBadRequest, "invalid username param")
		}

		queryLimit := 1
		if limit, err := strconv.ParseInt(c.QueryParam("limit"), 10, 8); err == nil && limit > 0 && limit <= 127 {
			queryLimit = int(limit)
		}

		tweets, err := nitter.NittosTweetsScrap(username, queryLimit)
		if err != nil {
			return res.HandleResp(http.StatusNotFound, "no tweets returned for this user")
		}

		res.SetPublicCache(15 * 60) // 15min
		return res.HandleResp(http.StatusOK, tweets)
	})

	n.GET("/nittos/:name/neets/:id", func(c echo.Context) error {
		res := routes.EchoWrapper{Context: c}

		username, neetId := c.Param("name"), c.Param("id")
		if utils.IsEmptyString(username) || !utils.IsSafeString(username) || utils.IsEmptyString(neetId) || len(neetId) < 19 {
			return res.HandleResp(http.StatusBadRequest, "invalid params")
		}

		queryLimit := 1
		if limit, err := strconv.ParseInt(c.QueryParam("limit"), 10, 8); err == nil && limit > 0 && limit <= 127 {
			queryLimit = int(limit)
		}

		comments, err := nitter.GetNeetComments(username, neetId, queryLimit)
		if err != nil {
			return res.HandleResp(http.StatusNotFound, "no comments returned for this neet")
		}

		res.SetPublicCache(30 * 60) // 30min
		return res.HandleResp(http.StatusOK, comments)
	})

	console.Log("NitterHandler Registered", console.Info)
}
