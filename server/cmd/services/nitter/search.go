package nitter

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"teniditter-server/cmd/global/utils"
	"teniditter-server/cmd/services"
	"teniditter-server/cmd/services/xml"

	"github.com/PuerkitoBio/goquery"
)

// q must be encoded before calling this func
func SearchTweets(q string) ([]TweetItem, error) {
	URL := fmt.Sprintf("https://nitter.pussthecat.org/search/rss?f=tweets&q=%s", q)
	if !utils.IsValidURL(URL) {
		return nil, errors.New("invalid URL")
	}

	resp, err := http.Get(URL)
	if err != nil || resp.StatusCode != 200 {
		return nil, err
	}
	defer resp.Body.Close()

	rawXml, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	tweetsXml, err := xml.ParseRSS[TweetItem](rawXml)
	if err != nil {
		return nil, err
	}

	for _, tweet := range tweetsXml.Channel.Items {
		if utils.ContainsScript(tweet.Desc) {
			return nil, errors.New("dangerous html detected")
		}
	}

	return tweetsXml.Channel.Items, nil
}

type NittosPreview struct {
	Username    string `json:"username"`
	Description string `json:"description"`
	AvatarUrl   string `json:"avatarUrl"`
}

func SearchNittos(username string) (*[]NittosPreview, error) {
	URL := fmt.Sprintf("https://nitter.pussthecat.org/search?f=users&q=%s", utils.FormatString(username))
	doc, err := services.GetHTMLDocument(URL)
	if err != nil {
		return nil, err
	}

	result := []NittosPreview{}
	doc.Find(".timeline-item").Each(func(i int, s *goquery.Selection) {
		username := s.Find(".tweet-body .tweet-header .username").Text()
		desc, _ := s.Find(".tweet-body .tweet-content").Html()
		if utils.ContainsScript(desc) {
			desc = ""
		}

		var avatarUrl string
		if avatarRaw, ok := s.Find(".tweet-body .tweet-header .tweet-avatar img").Attr("src"); ok {
			avatarUrl = "https://nitter.pussthecat.org" + avatarRaw
		}

		result = append(result, NittosPreview{username, desc, avatarUrl})
	})

	if len(result) <= 0 {
		return nil, errors.New("no users found")
	}

	return &result, nil
}
