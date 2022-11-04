package nitter

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"teniditter-server/cmd/global/utils"
	"teniditter-server/cmd/services/xml"
)

type Tweets xml.Rss[xml.TweetItem]

// q must be encoded before calling this func
func SearchTweets(q string) ([]xml.TweetItem, error) {
	if utils.IsEmptyString(q) {
		return nil, errors.New("invalid query param")
	}

	URL := fmt.Sprintf("https://nitter.pussthecat.org/search/rss?f=tweets&q=%s", q)
	resp, err := http.Get(URL)
	if err != nil || resp.StatusCode != 200 {
		return nil, err
	}
	defer resp.Body.Close()

	rawXml, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	searchData, err := xml.ParseNitterSearch(rawXml)
	if err != nil {
		return nil, err
	}
	return searchData.Channel.Items, nil
}
