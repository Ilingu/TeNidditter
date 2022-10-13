package teddit

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"teniditter-server/cmd/api/routes"
	"teniditter-server/cmd/global/console"
	"teniditter-server/cmd/global/utils"

	"github.com/PuerkitoBio/goquery"
	"github.com/labstack/echo/v4"
)

func TedditHandler(t *echo.Group) {
	t.GET("/r/:subreddit", func(c echo.Context) error {
		res := routes.EchoWrapper{Context: c}

		subreddit := utils.FormatUsername(c.Param("subreddit"))
		if utils.IsEmptyString(subreddit) {
			return res.HandleResp(http.StatusBadRequest, "invalid subreddit")
		}

		Url := fmt.Sprintf("https://teddit.net/r/%s", url.QueryEscape(subreddit))
		htmlPage, err := http.Get(Url)
		if err != nil || htmlPage.StatusCode != 200 {
			return res.HandleResp(http.StatusNotFound, err.Error())
		}
		defer htmlPage.Body.Close()

		doc, err := goquery.NewDocumentFromReader(htmlPage.Body)
		if err != nil {
			return res.HandleResp(500, err.Error())
		}

		subs := doc.Find("#sidebar .content > p:first-child").Text()
		description := doc.Find("#sidebar .content .heading").Text()
		rules, err := doc.Find("#sidebar .content .description").Html()
		if err != nil {
			return res.HandleResp(http.StatusNotFound, err.Error())
		}

		type subredditInfos struct {
			Subs        string `json:"subs"`
			Description string `json:"description"`
			Rules       string `json:"rules"`
		}

		respPayload := subredditInfos{
			Subs:        subs,
			Description: description,
			Rules:       rules,
		}
		return res.HandleResp(http.StatusOK, respPayload)
	})

	t.GET("/home", func(c echo.Context) error {
		res := routes.EchoWrapper{Context: c}

		homeFeedType := utils.SafeString(c.QueryParam("type"))
		afterId := utils.SafeString(c.QueryParam("afterId"))

		if !utils.IsEmptyString(homeFeedType) {
			valid := false
			for _, ftype := range []string{"hot", "new", "top", "rising", "controversial"} {
				if ftype == homeFeedType {
					valid = true
					break
				}
			}
			if !valid {
				return res.HandleResp(http.StatusBadRequest, "'type' query parameter is invalid")
			}
		} else {
			homeFeedType = "hot"
		}

		url := fmt.Sprintf("https://teddit.net/%s?api&raw_json=1", homeFeedType)
		if !utils.IsEmptyString(afterId) {
			url += fmt.Sprintf("&t=&after=t3_%s", afterId)
		}

		if !utils.IsValidURL(url) {
			return res.HandleResp(http.StatusForbidden, "This request has returned nothing")
		}

		resp, err := http.Get(url)
		if err != nil {
			return res.HandleResp(http.StatusForbidden, "This request has returned nothing")
		}
		defer resp.Body.Close()

		jsonBlob, err := io.ReadAll(resp.Body)
		if err != nil {
			return res.HandleResp(http.StatusForbidden, "This request has returned a corrupted response")
		}

		var jsonDatas map[string]interface{}
		err = json.Unmarshal(jsonBlob, &jsonDatas)
		if err != nil {
			return res.HandleResp(http.StatusForbidden, "This request has returned a corrupted response")
		}

		return res.HandleRespBlob(http.StatusOK, jsonDatas)
	})

	console.Log("TedditHandler Registered", console.Success)
}
