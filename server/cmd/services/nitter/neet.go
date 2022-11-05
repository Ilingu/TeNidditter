package nitter

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"teniditter-server/cmd/global/utils"
	"teniditter-server/cmd/services"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type NeetComments struct {
	MainThread []NeetComment   `json:"main"`
	Reply      [][]NeetComment `json:"reply"`
}

type NeetComment struct {
	NeetBasicComment
	Quote *NeetBasicComment `json:"quote,omitempty"`
}
type NeetBasicComment struct {
	Title       string           `json:"title"`
	Creator     NittosPreview    `json:"creator"`
	CreatedAt   int              `json:"createdAt"`
	Stats       NeetCommentStats `json:"stats"`
	Attachments *Attachments     `json:"attachment,omitempty"`
}
type Attachments struct {
	ImagesUrls []string `json:"images,omitempty"`
	VideosUrls []string `json:"videos,omitempty"`
}
type NeetCommentStats struct {
	ReplyCounts  int `json:"reply_counts"`
	RTCounts     int `json:"rt_counts"`
	QuotesCounts int `json:"quotes_counts,omitempty"`
	LikesCounts  int `json:"likes_counts"`
	PlayCounts   int `json:"play_counts,omitempty"`
}

func GetNeetContext(nittos, neetId string) (*NeetComment, error) {
	URL := fmt.Sprintf("https://nitter.pussthecat.org/%s/status/%s", nittos, neetId)
	doc, err := services.GetHTMLDocument(URL)
	if err != nil {
		return nil, err
	}

	commentSelector := doc.Find("#m > .timeline-item").First()
	commentData := extractCommentDatas(commentSelector)
	return &commentData, nil
}

func GetNeetComments(nittos, neetId string, limit int) (*NeetComments, error) {
	mainTheadSelector, repliesSelectors := queryCommentsSelectors(nittos, neetId, limit)
	if mainTheadSelector == nil || repliesSelectors == nil {
		return nil, errors.New("no selectors returned")
	}

	MainThread, Reply := make([]NeetComment, mainTheadSelector.Length()), make([][]NeetComment, repliesSelectors.Length())

	var wg sync.WaitGroup
	wg.Add(mainTheadSelector.Length() + repliesSelectors.Length())

	go mainTheadSelector.Each(func(i int, s *goquery.Selection) {
		go func() {
			defer wg.Done()
			MainThread[i] = extractCommentDatas(s)
		}()
	})

	go repliesSelectors.Each(func(i int, s *goquery.Selection) {
		go func() {
			defer wg.Done()

			ReplyGroup := []NeetComment{}
			s.Find(".timeline-item").Each(func(i int, t *goquery.Selection) {
				if !t.HasClass("more-replies") {
					ReplyGroup = append(ReplyGroup, extractCommentDatas(t))
				}
			})
			Reply[i] = ReplyGroup
		}()
	})

	wg.Wait()

	return &NeetComments{MainThread, Reply}, nil
}

// This function is error-less so check by yourself if the returned value are nil.
//
// It'll query the commentsSelector of a post by scrapping the page and hitting the "show more" button, and repeat this process "limit" times. It then return the main thread selector and the concatenated version of all the commentsSelectors
func queryCommentsSelectors(nittos, neetId string, limit int) (mainTheadSelector, commentsSelector *goquery.Selection) {
	var allComments *goquery.Selection

	URL := fmt.Sprintf("https://nitter.pussthecat.org/%s/status/%s", nittos, neetId)
	for i := 0; i < limit; i++ {
		doc, err := services.GetHTMLDocument(URL)
		if err != nil {
			return mainTheadSelector, allComments
		}

		replies := doc.Find("#r .reply")
		if allComments == nil {
			allComments = replies
			mainTheadSelector = doc.Find("div.main-thread .timeline-item")
		} else {
			allComments = allComments.AddNodes(replies.Nodes...)
		}
		nextQuery, exist := doc.Find("#r .show-more > a").Attr("href")
		if !exist {
			return mainTheadSelector, allComments // no comments left, return comments already fetched
		}

		nextUrl, err := url.Parse(URL)
		if err != nil {
			return mainTheadSelector, allComments
		}

		q := nextUrl.Query()
		q.Set("cursor", strings.TrimPrefix(nextQuery, "?cursor="))
		nextUrl.RawQuery, _ = url.QueryUnescape(q.Encode())

		URL = nextUrl.String()
	}

	return mainTheadSelector, allComments
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
	content, _ := s.Find(selector + "> .tweet-content").Html()
	if utils.ContainsScript(content) {
		content = ""
	}

	// Attachment
	imgUrl, vidUrl := []string{}, []string{}
	s.Find(selector + "> .attachments img").Each(func(i int, s *goquery.Selection) {
		if rawUrl, exist := s.Attr("src"); exist {
			imgUrl = append(imgUrl, "https://nitter.pussthecat.org"+rawUrl)
		}
	})
	s.Find(selector + "> .attachments video").Each(func(i int, s *goquery.Selection) {
		if rawUrl, exist := s.Attr("data-url"); exist {
			vidUrl = append(vidUrl, "https://nitter.pussthecat.org"+rawUrl)
		}
	})

	var attachments *Attachments
	if len(imgUrl) > 0 || len(vidUrl) > 0 {
		attachments = &Attachments{imgUrl, vidUrl}
	}

	// Stats
	var replyCounts, rtCounts, quotesCounts, likesCounts, playCounts int
	s.Find(selector + "> .tweet-stats .tweet-stat").Each(func(i int, s *goquery.Selection) {
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

	// Potential Quote
	var quote *NeetBasicComment
	if quoteUrl, exist := s.Find(selector + "> .quote > a.quote-link").Attr("href"); exist {
		if quoteData, err := fetchQuote(quoteUrl); err == nil {
			quote = &quoteData.NeetBasicComment
		}
	}

	commentData := NeetBasicComment{content, creator, int(createdAt), stats, attachments}
	return NeetComment{NeetBasicComment: commentData, Quote: quote}
}
