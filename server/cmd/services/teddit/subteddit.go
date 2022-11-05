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
	ps "teniditter-server/cmd/planetscale"
	"teniditter-server/cmd/redis"
	"teniditter-server/cmd/redis/rediskeys"
	"teniditter-server/cmd/services"
	"time"
)

func GetSubredditPosts(subreddit string) (*map[string]any, error) {
	redisKey := rediskeys.NewKey(rediskeys.SUBREDDIT, subreddit+"_POSTS")

	if posts, err := redis.Get[map[string]any](redisKey); err == nil {
		console.Log("Subteddit Posts Returned from Cache ⚡", console.Neutral)
		return &posts, nil
	}

	Url := fmt.Sprintf("https://teddit.net/r/%s?api&raw_json=1", url.QueryEscape(subreddit))
	if !utils.IsValidURL(Url) {
		return nil, errors.New("invalid URL")
	}

	rawPosts, err := http.Get(Url)
	if err != nil || rawPosts.StatusCode != 200 {
		return nil, err
	}
	defer rawPosts.Body.Close()

	rawBlobPosts, err := io.ReadAll(rawPosts.Body)
	if err != nil {
		return nil, err
	}

	var posts map[string]any
	err = json.Unmarshal(rawBlobPosts, &posts)
	if err != nil || len(posts) <= 0 {
		return nil, err
	}

	// Caching
	go redis.Set(redisKey, posts, 24*time.Hour)

	return &posts, nil
}

type subredditInfos struct {
	Subs        string `json:"subs"`
	Description string `json:"description"`
	Rules       string `json:"rules"`
}

func GetSubredditMetadatas(subreddit string) (*subredditInfos, error) {
	redisKey := rediskeys.NewKey(rediskeys.SUBREDDIT, subreddit+"_ABOUT")

	if subDatas, err := redis.Get[subredditInfos](redisKey); err == nil {
		console.Log("Subteddit Returned from Cache ⚡", console.Neutral)
		return &subDatas, nil
	}

	URL := fmt.Sprintf("https://teddit.net/r/%s", url.QueryEscape(subreddit))
	doc, err := services.GetHTMLDocument(URL)
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
	if utils.ContainsScript(rules) {
		return nil, errors.New("dangerous html detected for rules")
	}

	respPayload := subredditInfos{
		Subs:        subs,
		Description: description,
		Rules:       rules,
	}

	// Caching
	go redis.Set(redisKey, respPayload, 7*24*time.Hour) // 7d
	go func(subname string) {
		db := ps.DBManager.Connect()
		if db != nil {
			db.Exec("INSERT INTO Subteddits (subname) VALUES (?);", subname)
		}
	}(subreddit)

	return &respPayload, nil
}
