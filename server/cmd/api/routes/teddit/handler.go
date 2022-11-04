package teddit_routes

import (
	"net/http"
	"teniditter-server/cmd/api/routes"
	"teniditter-server/cmd/db"
	"teniditter-server/cmd/global/console"
	"teniditter-server/cmd/global/utils"
	"teniditter-server/cmd/services/teddit"

	"github.com/labstack/echo/v4"
)

func TedditHandler(t *echo.Group) {
	t.GET("/u/:username", func(c echo.Context) error {
		res := routes.EchoWrapper{Context: c}

		username := c.Param("username")
		if utils.IsEmptyString(username) {
			return res.HandleResp(http.StatusBadRequest, "invalid username")
		}

		userInfos, err := teddit.GetUserInfos(username)
		if err != nil {
			return res.HandleResp(http.StatusForbidden, err.Error())
		}

		res.SetPublicCache(48 * 60 * 60) // 48h
		return res.HandleResp(http.StatusOK, userInfos)
	})

	t.GET("/r/:subreddit/post/:id", func(c echo.Context) error {
		res := routes.EchoWrapper{Context: c}

		subreddit := utils.FormatToSafeString(c.Param("subreddit"))
		postId := utils.SafeString(c.Param("id"))
		sort := c.QueryParam("sort")

		if utils.IsEmptyString(subreddit) || utils.IsEmptyString(postId) || len(postId) < 6 {
			return res.HandleResp(http.StatusBadRequest, "invalid subreddit or postId")
		}

		if sort != "best" && sort != "top" && sort != "new" && sort != "controversial" && sort != "old" && sort != "qa" {
			sort = "best"
		}

		postDatas, err := teddit.GetPostInfo(subreddit, postId, sort)
		if err != nil {
			return res.HandleResp(http.StatusBadRequest, "no data returned")
		}

		res.SetPublicCache(48 * 60 * 60) // 48h
		return res.HandleResp(http.StatusOK, postDatas)
	})

	t.GET("/r/search", func(c echo.Context) error {
		res := routes.EchoWrapper{Context: c}

		subreddit := utils.FormatToSafeString(c.QueryParam("q"))
		if utils.IsEmptyString(subreddit) {
			return res.HandleResp(http.StatusBadRequest, "invalid subreddit")
		}

		matchSubs, err := db.SearchSubteddit(subreddit)
		if err != nil {
			return res.HandleResp(http.StatusForbidden, err.Error())
		}

		res.SetPublicCache(30 * 60) // 30min
		return res.HandleResp(http.StatusOK, matchSubs)
	})

	t.GET("/r/:subreddit/about", func(c echo.Context) error {
		res := routes.EchoWrapper{Context: c}

		subreddit := utils.FormatToSafeString(c.Param("subreddit"))
		if utils.IsEmptyString(subreddit) {
			return res.HandleResp(http.StatusBadRequest, "invalid subreddit")
		}

		subredditMetadatas, err := teddit.GetSubredditMetadatas(subreddit)
		if err != nil {
			return res.HandleResp(http.StatusForbidden, err.Error())
		}

		res.SetPublicCache(7 * 24 * 60 * 60) // 1w
		return res.HandleResp(http.StatusOK, subredditMetadatas)
	})

	t.GET("/r/:subreddit/posts", func(c echo.Context) error {
		res := routes.EchoWrapper{Context: c}

		subreddit := utils.FormatToSafeString(c.Param("subreddit"))
		if utils.IsEmptyString(subreddit) {
			return res.HandleResp(http.StatusBadRequest, "invalid subreddit")
		}

		subredditPosts, err := teddit.GetSubredditPosts(subreddit)
		if err != nil {
			return res.HandleResp(http.StatusForbidden, err.Error())
		}

		res.SetPublicCache(12 * 60 * 60) // 1/2d
		return res.HandleRespBlob(http.StatusOK, *subredditPosts)
	})

	t.GET("/home", func(c echo.Context) error {
		res := routes.EchoWrapper{Context: c}

		FeedType := "hot"
		afterId := utils.SafeString(c.QueryParam("afterId"))

		if ft := utils.SafeString(c.QueryParam("type")); !utils.IsEmptyString(ft) {
			if ft != "hot" && ft != "new" && ft != "top" && ft != "rising" && ft != "controversial" {
				return res.HandleResp(http.StatusBadRequest, "'type' query parameter is invalid")
			}
			FeedType = ft
		}

		posts, err := teddit.GetHomePosts(FeedType, afterId, false)
		if err != nil {
			return res.HandleResp(http.StatusForbidden, err.Error())
		}

		res.SetPublicCache(1 * 60 * 60) // 1h
		return res.HandleRespBlob(http.StatusOK, *posts)
	})

	console.Log("TedditHandler Registered", console.Info)
}
