package nitter_routes

import (
	"errors"
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
			return res.HandleResp(http.StatusNotImplemented)
		default:
			return res.HandleResp(http.StatusBadRequest, `invalid query "type", it must be whether "tweets" or "users"`)
		}
	})

	n.GET("/nittos/:name", func(c echo.Context) error { return errors.New("") })
	n.GET("/nittos/:name/:neetId", func(c echo.Context) error { return errors.New("") })

	console.Log("NitterHandler Registered", console.Info)
}
