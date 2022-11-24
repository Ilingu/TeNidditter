package nitter

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"teniditter-server/cmd/global/console"
	"teniditter-server/cmd/global/utils"
	"teniditter-server/cmd/redis"
	"teniditter-server/cmd/redis/rediskeys"
	"teniditter-server/cmd/services"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func GetNeetContext(nittos, neetId string) (*NeetComment, error) {
	URL := fmt.Sprintf("https://nitter.pussthecat.org/%s/status/%s", nittos, neetId)
	doc, err := services.GetHTMLDocument(URL)
	if err != nil {
		return nil, err
	}

	commentSelector := doc.Find("#m > .timeline-item").First()

	commentData := extractNeetDatas(commentSelector)
	commentData.Id = neetId

	return &commentData, nil
}

func GetNeetComments(nittos, neetId string, limit int) (*NeetInfo, error) {
	redisKey := rediskeys.NewKey(rediskeys.NITTER_NEET_COMMENTS, utils.GenerateKeyFromArgs(nittos, neetId, limit))
	if neetInfo, err := redis.Get[NeetInfo](redisKey); err == nil {
		console.Log("Neet Returned from cache", console.Neutral)
		return &neetInfo, nil // Returned from cache
	}

	URL := fmt.Sprintf("https://nitter.pussthecat.org/%s/status/%s", nittos, neetId)
	replies, nextQuery := "#r .reply", "#r .show-more > a"

	doc, repliesSelectors := queryMoreSelectors(URL, replies, nextQuery, limit)
	if doc == nil {
		return nil, errors.New("no doc returned")
	}

	mainTheadSelector := doc.Find("div.main-thread .timeline-item")
	if mainTheadSelector == nil || repliesSelectors == nil {
		return nil, errors.New("no selectors returned")
	}
	if mainTheadSelector.Length() <= 0 {
		return nil, errors.New("context tweets not found")
	}

	MainThread, Reply := make([]NeetComment, mainTheadSelector.Length()), make([][]NeetComment, repliesSelectors.Length())

	var wg sync.WaitGroup
	wg.Add(mainTheadSelector.Length() + repliesSelectors.Length())

	go mainTheadSelector.Each(func(i int, s *goquery.Selection) {
		go func() {
			defer wg.Done()

			neetDatas := extractNeetDatas(s)
			neetDatas.Id = neetId
			MainThread[i] = neetDatas
		}()
	})

	go repliesSelectors.Each(func(i int, s *goquery.Selection) {
		go func() {
			defer wg.Done()

			ReplyGroup := []NeetComment{}
			s.Find(".timeline-item").Each(func(i int, t *goquery.Selection) {
				if !t.HasClass("more-replies") {
					ReplyGroup = append(ReplyGroup, extractNeetDatas(t))
				}
			})
			Reply[i] = ReplyGroup
		}()
	})

	wg.Wait()

	deleteDuplicatesNeets(&Reply)
	result := NeetInfo{MainThread, Reply}

	// caching
	exp := time.Hour
	if MainThread[0].Stats.LikesCounts < 100 {
		exp = 12 * time.Hour
	}
	go redis.Set(redisKey, result, exp)

	return &result, nil
}

func fetchNeetThread(s *goquery.Selection) ([]NeetComment, int) {
	currComment := extractNeetDatas(s)
	thread, toExclude := []NeetComment{currComment}, 0
	if s.HasClass("thread") && !s.HasClass("thread-last") && s.Next().Length() != 0 {
		toExclude++
		childThread, childToExclude := fetchNeetThread(s.Next())

		thread = append(thread, childThread...)
		toExclude += childToExclude
	}

	return thread, toExclude
}

func extractNeetDatas(s *goquery.Selection) NeetComment {
	// ID
	var id string
	idLink := s.Find("a.tweet-link").AttrOr("href", "")
	if paths := strings.Split(idLink, "/"); len(paths) > 0 {
		id = strings.TrimSuffix(paths[len(paths)-1], "#m")
	}

	selector := ".tweet-body "

	// Header (creator, createdAt)
	nittos := s.Find(selector + "> div a.username").First().Text()
	pinned := s.Find(selector+"> div > .pinned").Length() == 1

	var retweetedBy string
	if retweet := utils.TrimString(s.Find(selector + "> div > .retweet-header").First().Text()); !utils.IsEmptyString(retweet) {
		retweetedBy = strings.TrimSuffix(retweet, " retweeted")
	}

	var avatarUrl string
	if avatarRaw, ok := s.Find(selector + "> div a.tweet-avatar > img.avatar").First().Attr("src"); ok {
		avatarUrl = "https://nitter.pussthecat.org" + avatarRaw
	}
	creator := NittosPreview{Username: nittos, AvatarUrl: avatarUrl}

	var createdAt int64
	if dateFormatted, exist := s.Find(selector + "> div span.tweet-date > a").First().Attr("title"); exist {
		const layout = "Jan 2, 2006 · 3:04 PM UTC"
		if t, err := time.Parse(layout, dateFormatted); err == nil {
			createdAt = t.Unix()
		}
	}

	// Body/Title
	content, _ := s.Find(selector + "> .tweet-content").First().Html()
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
			if paths := strings.Split(rawUrl, "/"); len(paths) > 0 {
				if hostUrl, err := url.QueryUnescape(paths[len(paths)-1]); err == nil {
					vidUrl = append(vidUrl, hostUrl)
				}
			}
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
		case 0:
			replyCounts = num
		case 1:
			rtCounts = num
		case 2:
			quotesCounts = num
		case 3:
			likesCounts = num
		case 4:
			playCounts = num
		}
	})
	stats := NeetCommentStats{replyCounts, rtCounts, quotesCounts, likesCounts, playCounts}

	// Potential Quote
	var quote *NeetBasicComment
	if quoteUrl, exist := s.Find(selector + "> .quote > a.quote-link").First().Attr("href"); exist {
		if quoteData, err := fetchCtxNeetFromUrl(quoteUrl); err == nil {
			quote = &quoteData.NeetBasicComment
		}
	}

	// Potential Link Card (external link)
	linkCard, _ := s.Find(selector + "> .card > a.card-container").First().Attr("href")

	commentData := NeetBasicComment{id, content, creator, int(createdAt), stats, attachments, linkCard, retweetedBy, pinned}
	return NeetComment{NeetBasicComment: commentData, Quote: quote}
}
