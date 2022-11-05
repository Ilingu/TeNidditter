package nitter_routes

import (
	"net/http"
	"net/url"
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

		if !utils.IsUrlEncoded(QuerySearch) {
			QuerySearch = url.QueryEscape(QuerySearch)
		}

		queryLimit := 1
		if limit, err := strconv.ParseInt(c.QueryParam("limit"), 10, 8); err == nil && limit > 0 && limit <= 127 {
			queryLimit = int(limit)
		}

		switch SearchType {
		case "tweets":
			var tweets []nitter.NeetComment
			var err error
			if queryLimit > 1 {
				res.SetPublicCache(15 * 60) // 15min
				tweets, err = nitter.SearchTweetsScrap(QuerySearch, queryLimit)
			} else {
				res.SetPublicCache(1 * 60 * 60) // 1h
				tweets, err = nitter.SearchTweetsXML(QuerySearch)
			}
			if err != nil {
				res.Response().Header().Del("Cache-Control")
				return res.HandleResp(http.StatusNotFound, err.Error())
			}

			return res.HandleResp(http.StatusOK, tweets)
		case "users":
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

		username := utils.FormatString(c.Param("name"))
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

		username := utils.FormatString(c.Param("name"))
		if utils.IsEmptyString(username) {
			return res.HandleResp(http.StatusBadRequest, "invalid username param")
		}

		queryLimit := 1
		if limit, err := strconv.ParseInt(c.QueryParam("limit"), 10, 8); err == nil && limit > 0 && limit <= 127 {
			queryLimit = int(limit)
		}

		var tweets []nitter.NeetComment
		var err error
		if queryLimit > 1 {
			res.SetPublicCache(15 * 60) // 15min
			tweets, err = nitter.NittosTweetsScrap(username, queryLimit)
		} else {
			res.SetPublicCache(1 * 60 * 60) // 1h
			tweets, err = nitter.NittosTweetsXML(username)
		}
		if err != nil {
			res.Response().Header().Del("Cache-Control")
			return res.HandleResp(http.StatusNotFound, "no tweets returned for this user")
		}

		return res.HandleResp(http.StatusOK, tweets)
	})

	n.GET("/nittos/:name/neets/:id", func(c echo.Context) error {
		res := routes.EchoWrapper{Context: c}

		username, neetId := utils.FormatString(c.Param("name")), c.Param("id")
		if utils.IsEmptyString(username) || utils.IsEmptyString(neetId) || len(neetId) < 19 {
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
