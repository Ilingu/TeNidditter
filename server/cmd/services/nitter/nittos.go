package nitter

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"teniditter-server/cmd/global/console"
	"teniditter-server/cmd/global/utils"
	utils_enc "teniditter-server/cmd/global/utils/encryption"
	ps "teniditter-server/cmd/planetscale"
	"teniditter-server/cmd/redis"
	"teniditter-server/cmd/redis/rediskeys"
	"teniditter-server/cmd/services"
	"teniditter-server/cmd/services/xml"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func NittosMetadata(username string) (*Nittos, error) {
	URL := fmt.Sprintf("https://nitter.pussthecat.org/%s", username)
	doc, err := services.GetHTMLDocument(URL)
	if err != nil {
		return nil, err
	}

	nittosName := doc.Find("div.profile-tab.sticky .profile-card a.profile-card-username").Text()

	var bannerUrl string
	if bannerRaw, ok := doc.Find("div.profile-banner > a").Attr("href"); ok {
		bannerUrl = "https://nitter.pussthecat.org" + bannerRaw
	}

	var avatarUrl string
	if avatarRaw, ok := doc.Find("div.profile-tab.sticky .profile-card a.profile-card-avatar").Attr("href"); ok {
		avatarUrl = "https://nitter.pussthecat.org" + avatarRaw
	}

	metadataSelector := "div.profile-tab.sticky .profile-card .profile-card-extra "

	bio, _ := doc.Find(metadataSelector + ".profile-bio").Html()
	if utils.ContainsScript(bio) {
		bio = ""
	}

	location := utils.TrimString(doc.Find(metadataSelector + ".profile-location").Text())
	websiteLink, _ := doc.Find(metadataSelector + ".profile-website a").Attr("href")
	joinDate := utils.TrimString(doc.Find(metadataSelector + ".profile-joindate").Text())

	var tweetsCount, followingCounts, followersCounts, likesCounts int
	StatSelector := doc.Find(metadataSelector + ".profile-card-extra-links > ul .profile-stat-num")
	StatSelector.Each(func(i int, s *goquery.Selection) {
		num, _ := strconv.Atoi(strings.ReplaceAll(utils.TrimString(s.Text()), ",", ""))
		switch i {
		case 0:
			tweetsCount = num
		case 1:
			followingCounts = num
		case 2:
			followersCounts = num
		case 3:
			likesCounts = num
		}
	})

	stats := NittosStats{tweetsCount, followingCounts, followersCounts, likesCounts}
	metadata := Nittos{nittosName, bio, avatarUrl, location, websiteLink, joinDate, stats, bannerUrl}

	// Caching
	go func() {
		db := ps.DBManager.Connect()
		if db != nil {
			db.Exec("INSERT INTO Twittos (username) VALUES (?);", username)
		}
	}()

	return &metadata, nil
}

func NittosTweetsScrap(username string, limit int) ([][]NeetComment, error) {
	redisKey := rediskeys.NewKey(rediskeys.NITTER_NITTOS_TWEETS, utils_enc.GenerateHashFromArgs(username, limit))
	if comments, err := redis.Get[[][]NeetComment](redisKey); err == nil {
		console.Log("Neets Returned from cache", console.Neutral)
		return comments, nil // Returned from cache
	}

	URL := fmt.Sprintf("https://nitter.pussthecat.org/%s", username)
	tweets, err := fetchTweets(URL, limit)
	if err != nil {
		return nil, err
	}

	// Caching
	go redis.Set(redisKey, tweets, 30*time.Minute)

	return tweets, nil
}

func NittosTweetsXML(username string) ([]NeetComment, error) {
	URL := fmt.Sprintf("https://nitter.pussthecat.org/%s/rss", username)
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
