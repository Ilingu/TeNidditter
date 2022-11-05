package nitter

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"teniditter-server/cmd/global/utils"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type NeetComments struct {
	MainThread []NeetComment `json:"main"`
	Reply      []NeetComment `json:"reply"`
}

type NeetComment struct {
	Title       string           `json:"title"`
	Creator     NittosPreview    `json:"creator"`
	CreatedAt   int              `json:"createdAt"`
	Stats       NeetCommentStats `json:"stats"`
	Attachments Attachments      `json:"attachment"`
}
type Attachments struct {
	ImagesUrls []string `json:"images"`
	VideosUrls []string `json:"videos"`
}
type NeetCommentStats struct {
	ReplyCounts  int `json:"reply_counts"`
	RTCounts     int `json:"rt_counts"`
	QuotesCounts int `json:"quotes_counts"`
	LikesCounts  int `json:"likes_counts"`
	PlayCounts   int `json:"play_counts"`
}

func GetNeetComments(nittos, neetId string) (*NeetComments, error) {
	URL := fmt.Sprintf("https://nitter.pussthecat.org/%s/status/%s#m", nittos, neetId)
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

	mainTheadSelector := doc.Find("div.main-thread .timeline-item")
	replySelector := doc.Find("#r .reply")

	MainThread, Reply := []NeetComment{}, []NeetComment{}

	var wg sync.WaitGroup
	wg.Add(mainTheadSelector.Length() + replySelector.Length())

	go mainTheadSelector.Each(func(i int, s *goquery.Selection) {
		go func() {
			defer wg.Done()
			MainThread = append(MainThread, extractCommentDatas(s))
		}()
	})

	go replySelector.Each(func(i int, s *goquery.Selection) {
		go func() {
			defer wg.Done()
			Reply = append(Reply, extractCommentDatas(s))
		}()
	})

	wg.Wait()

	return &NeetComments{MainThread, Reply}, nil
}

func extractCommentDatas(s *goquery.Selection) NeetComment {
	selector := ".tweet-body "

	// Header (creator, createdAt)
	nittos := s.Find(selector + "> div a.username").Text()

	var avatarUrl string
	if avatarRaw, ok := s.Find(selector + "> div a.tweet-avatar > img.avatar").Attr("src"); ok {
		avatarUrl = "https://nitter.pussthecat.org" + avatarRaw
	}
	creator := NittosPreview{Username: nittos, AvatarUrl: avatarUrl}

	var createdAt int64
	if dateFormatted, exist := s.Find(selector + "> div span.tweet-date > a").Attr("title"); exist {
		const layout = "Jan 2, 2006 · 3:04 PM UTC"
		if t, err := time.Parse(layout, dateFormatted); err == nil {
			createdAt = t.Unix()
		}
	}

	// Body/Title
	content, _ := s.Find(selector + ".tweet-content").Html()
	if utils.ContainsScript(content) {
		content = ""
	}

	// Attachment
	imgUrl, vidUrl := []string{}, []string{}
	s.Find(selector + ".attachments img").Each(func(i int, s *goquery.Selection) {
		if rawUrl, exist := s.Attr("src"); exist {
			imgUrl = append(imgUrl, "https://nitter.pussthecat.org"+rawUrl)
		}
	})
	s.Find(selector + ".attachments video").Each(func(i int, s *goquery.Selection) {
		if rawUrl, exist := s.Attr("data-url"); exist {
			vidUrl = append(vidUrl, "https://nitter.pussthecat.org"+rawUrl)
		}
	})
	attachments := Attachments{imgUrl, vidUrl}

	// Stats
	var replyCounts, rtCounts, quotesCounts, likesCounts, playCounts int
	s.Find(selector + ".tweet-stats .tweet-stat").Each(func(i int, s *goquery.Selection) {
		num, _ := strconv.Atoi(strings.ReplaceAll(utils.TrimString(s.Text()), ",", ""))
		switch i {
		case 1:
			replyCounts = num
		case 2:
			rtCounts = num
		case 3:
			quotesCounts = num
		case 4:
			likesCounts = num
		case 5:
			playCounts = num
		}
	})
	stats := NeetCommentStats{replyCounts, rtCounts, quotesCounts, likesCounts, playCounts}

	return NeetComment{content, creator, int(createdAt), stats, attachments}
}
