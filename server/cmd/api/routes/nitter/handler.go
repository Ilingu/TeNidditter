package nitter_routes

import (
	"net/http"
	"net/url"
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

		switch SearchType {
		case "tweets":
			tweets, err := nitter.SearchTweets(QuerySearch)
			if err != nil {
				return res.HandleResp(http.StatusNotFound, err.Error())
			}

			res.SetPublicCache(1 * 60 * 60) // 1h
			return res.HandleResp(http.StatusOK, tweets)
		case "users":
			nittos, err := nitter.SearchNittos(QuerySearch)
			if err != nil {
				return res.HandleResp(http.StatusNotFound, err.Error())
			}

			res.SetPublicCache(1 * 60 * 60) // 1h
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
		return res.HandleResp(http.StatusOK, metadata)
	})
	n.GET("/nittos/:name/neets", func(c echo.Context) error {
		res := routes.EchoWrapper{Context: c}

		username := utils.FormatString(c.Param("name"))
		if utils.IsEmptyString(username) {
			return res.HandleResp(http.StatusBadRequest, "invalid username param")
		}

		tweets, err := nitter.NittosTweets(username)
		if err != nil {
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

		comments, err := nitter.GetNeetComments(username, neetId)
		if err != nil {
			return res.HandleResp(http.StatusNotFound, "no comments returned for this neet")
		}

		return res.HandleResp(http.StatusOK, comments)
	})

	console.Log("NitterHandler Registered", console.Info)
}
