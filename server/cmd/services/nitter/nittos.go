package nitter

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"teniditter-server/cmd/global/utils"
	"teniditter-server/cmd/services/xml"

	"github.com/PuerkitoBio/goquery"
)

type Nittos struct {
	Username  string      `json:"username"`
	Bio       string      `json:"bio"`
	AvatarUrl string      `json:"avatarUrl"`
	Location  string      `json:"location"`
	Website   string      `json:"website"`
	JoinDate  string      `json:"joinDate"`
	Stats     NittosStats `json:"stats"`
	BannerUrl string      `json:"bannerUrl"`
}

type NittosStats struct {
	TweetsCounts    int `json:"tweets_counts"`
	FollowingCounts int `json:"following_counts"`
	FollowersCounts int `json:"followers_counts"`
	LikesCounts     int `json:"likes_counts"`
}

func NittosMetadata(username string) (*Nittos, error) {
	URL := fmt.Sprintf("https://nitter.pussthecat.org/%s", username)
	if !utils.IsValidURL(URL) {
		return nil, errors.New("invalid URL")
	}

	htmlPage, err := http.Get(URL)
	if err != nil || htmlPage.StatusCode != 200 {
		return nil, err
	}
	defer htmlPage.Body.Close()

	doc, err := goquery.NewDocumentFromReader(htmlPage.Body)
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
		case 1:
			tweetsCount = num
		case 2:
			followingCounts = num
		case 3:
			followersCounts = num
		case 4:
			likesCounts = num
		}
	})

	stats := NittosStats{tweetsCount, followingCounts, followersCounts, likesCounts}
	metadata := Nittos{nittosName, bio, avatarUrl, location, websiteLink, joinDate, stats, bannerUrl}

	return &metadata, nil
}

func NittosTweets(username string) ([]TweetItem, error) {
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
