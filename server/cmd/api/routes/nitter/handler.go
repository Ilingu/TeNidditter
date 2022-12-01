package nitter_routes

import (
	"net/http"
	"strconv"
	"sync"
	"teniditter-server/cmd/api/routes"
	"teniditter-server/cmd/api/sse"
	"teniditter-server/cmd/global/console"
	"teniditter-server/cmd/global/utils"
	"teniditter-server/cmd/services/nitter"
	"time"

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

			if clientIp, err := routes.GetIP(c.Request(), true); err == nil {
				go StreamExternalLinks(clientIp, tweets)
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

		if clientIp, err := routes.GetIP(c.Request(), true); err == nil {
			go StreamExternalLinks(clientIp, tweets)
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

		if clientIp, err := routes.GetIP(c.Request(), true); err == nil {
			go StreamExternalLinks(clientIp, [][]nitter.NeetComment{comments.MainThread})
			go StreamExternalLinks(clientIp, comments.Reply)
		}

		res.SetPublicCache(30 * 60) // 30min
		return res.HandleResp(http.StatusOK, comments)
	})

	n.GET("/stream-in-external-links", echo.WrapHandler(http.HandlerFunc(sse.SSEHandler)))
	console.Log("NitterHandler Registered")
}

// This will compute the external links datas (meta tags in html response)
//
// And then streams it back to the client thanks to [Server-Side Events] (only if client and server already/or will have a SSE conn initialized (10s max))
//
// [Server-Side Events]: https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events
func StreamExternalLinks(reqIp string, tweets [][]nitter.NeetComment) {
	client, exist := sse.GetClient(reqIp)
	if !exist {
		waitClientToConnect := make(chan bool)
		go sse.ListenToNewConn(func(ip string) {
			if ip == reqIp {
				client, exist = sse.GetClient(reqIp)
				waitClientToConnect <- exist
			}
		})

		timeout := time.After(20 * time.Second)
		select {
		case good := <-waitClientToConnect:
			if !good {
				return
			}
		case <-timeout:
			return
		}
	}

	var totalLength uint64
	for i := 0; i < len(tweets); i++ {
		totalLength += uint64(len(tweets[i]))
	}

	allExternalLinks := []*nitter.NeetBasicComment{}
	for _, t := range tweets {
		for _, n := range t {
			if !utils.IsEmptyString(n.ExternalLink) {
				allExternalLinks = append(allExternalLinks, &n.NeetBasicComment)
			}
			if n.Quote != nil && !utils.IsEmptyString(n.Quote.ExternalLink) {
				allExternalLinks = append(allExternalLinks, n.Quote)
			}
		}
	}

	if len(allExternalLinks) < 1 {
		client.Close <- true
		return
	}
	client.MinOpsNumber <- uint64(len(allExternalLinks))

	var wg sync.WaitGroup
	wg.Add(len(allExternalLinks))
	for _, n := range allExternalLinks {
		go func(neet *nitter.NeetBasicComment) {
			defer wg.Done()
			if metatagsDatas, err := nitter.GetExternalLinksMetatags(neet.ExternalLink); err == nil {
				metatagsDatas["neetId"] = neet.Id
				client.Events <- metatagsDatas
			}
		}(n)
	}

	wg.Wait()
	client.Close <- true
}
