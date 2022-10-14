package teddit

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"teniditter-server/cmd/global/console"
	"teniditter-server/cmd/global/utils"
	"teniditter-server/cmd/redis"
	"teniditter-server/cmd/redis/rediskeys"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func GetHomePosts(FeedType, afterId string) (*map[string]any, error) {
	// Check If content already cached:
	redisKey := rediskeys.NewKey(rediskeys.TEDDIT_HOME, FeedType)
	if posts, err := redis.Get[map[string]any](redisKey); err == nil {
		console.Log("Posts Returned from cache", console.Neutral)
		return &posts, nil // Returned from cache
	}

	url := fmt.Sprintf("https://teddit.net/%s?api&raw_json=1", FeedType)
	if !utils.IsEmptyString(afterId) {
		url += fmt.Sprintf("&t=&after=t3_%s", afterId)
	}

	if !utils.IsValidURL(url) {
		return nil, errors.New("this request has returned nothing")
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, errors.New("this request has returned nothing")
	}
	defer resp.Body.Close()

	jsonBlob, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("this request has returned a corrupted response")
	}

	var jsonDatas map[string]any
	err = json.Unmarshal(jsonBlob, &jsonDatas)
	if err != nil {
		return nil, errors.New("this request has returned a corrupted response")
	}

	// Caching
	go redis.Set(redisKey, jsonDatas)

	return &jsonDatas, nil
}

type subredditInfos struct {
	Subs        string `json:"subs"`
	Description string `json:"description"`
	Rules       string `json:"rules"`
}

func GetSubredditMetadatas(subreddit string) (*subredditInfos, error) {
	redisKey := rediskeys.NewKey(rediskeys.SUBREDDIT_INFO, subreddit)

	if subDatas, err := redis.Get[subredditInfos](redisKey); err == nil {
		console.Log("Subteddit Returned from Cache ⚡", console.Neutral)
		return &subDatas, nil
	}

	Url := fmt.Sprintf("https://teddit.net/r/%s", url.QueryEscape(subreddit))
	htmlPage, err := http.Get(Url)
	if err != nil || htmlPage.StatusCode != 200 {
		return nil, err
	}
	defer htmlPage.Body.Close()

	doc, err := goquery.NewDocumentFromReader(htmlPage.Body)
	if err != nil {
		return nil, err
	}

	if e := doc.Find(".reddit-error").Length(); e > 0 {
		return nil, errors.New("subreddit not found")
	}

	subs := doc.Find("#sidebar .content > p:first-child").Text()
	description := doc.Find("#sidebar .content .heading").Text()
	rules, err := doc.Find("#sidebar .content .description").Html()
	if err != nil {
		return nil, err
	}

	respPayload := subredditInfos{
		Subs:        subs,
		Description: description,
		Rules:       rules,
	}

	// Caching
	go redis.Set(redisKey, respPayload, 24*time.Hour)

	return &respPayload, nil
}
