package nitter

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"teniditter-server/cmd/global/console"
	"teniditter-server/cmd/global/utils"
	"teniditter-server/cmd/redis"
	"teniditter-server/cmd/redis/rediskeys"
	"teniditter-server/cmd/services/xml"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// q must be encoded before calling this func
func SearchTweetsXML(q string) ([]NeetComment, error) {
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

	tweetsXml, err := xml.ParseRSS[XmlTweetItem](rawXml)
	if err != nil {
		return nil, err
	}

	tweets := make([]NeetComment, len(tweetsXml.Channel.Items))
	for i, tweetXml := range tweetsXml.Channel.Items {
		if utils.ContainsScript(tweetXml.Desc) {
			return nil, errors.New("dangerous html detected")
		}

		tweet, _ := tweetXml.ToJSON(true) // no error possible when ToJSON is set to `true`
		tweets[i] = *tweet
	}

	return tweets, nil
}

func SearchTweetsScrap(q string, limit int) ([][]NeetComment, error) {
	redisKey := rediskeys.NewKey(rediskeys.NITTER_SEARCH_TWEETS, utils.GenerateKeyFromArgs(q, limit))
	if comments, err := redis.Get[[][]NeetComment](redisKey); err == nil {
		console.Log("Neets Returned from cache", console.Neutral)
		return comments, nil // Returned from cache
	}

	URL := fmt.Sprintf("https://nitter.pussthecat.org/search?f=tweets&q=%s", url.QueryEscape(q))
	tweets, err := fetchTweets(URL, limit)
	if err != nil {
		return nil, err
	}

	// Caching
	go redis.Set(redisKey, tweets, 30*time.Minute)

	return tweets, nil
}

func SearchNittos(username string, limit int) (*[]NittosPreview, error) {
	redisKey := rediskeys.NewKey(rediskeys.NITTER_SEARCH_NITTOS, utils.GenerateKeyFromArgs(username, limit))
	if nittos, err := redis.Get[[]NittosPreview](redisKey); err == nil {
		console.Log("Nittos Returned from cache", console.Neutral)
		return &nittos, nil // Returned from cache
	}

	URL := fmt.Sprintf("https://nitter.pussthecat.org/search?f=users&q=%s", username)

	_, nittosSelectors := queryMoreSelectors(URL, ".timeline-item", "div.timeline > .show-more:not(.timeline-item) > a", limit)
	if nittosSelectors == nil {
		return nil, errors.New("no tweets found")
	}

	result := []NittosPreview{}
	nittosSelectors.Each(func(i int, s *goquery.Selection) {
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

	// Caching
	go redis.Set(redisKey, result, 30*time.Minute)

	return &result, nil
}
