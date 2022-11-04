package teddit

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"sync"
	"teniditter-server/cmd/global/console"
	"teniditter-server/cmd/global/utils"
	"teniditter-server/cmd/redis"
	"teniditter-server/cmd/redis/rediskeys"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// FeedType is whether "hot" or "new" or "top" or "rising" or "controversial"
func GetHomePosts(FeedType, afterId string, nocache bool) (*map[string]any, error) {
	// Check If content already cached:
	redisKey := rediskeys.NewKey(rediskeys.TEDDIT_HOME, FeedType+afterId)
	if !nocache {
		if posts, err := redis.Get[map[string]any](redisKey); err == nil {
			console.Log("Posts Returned from cache", console.Neutral)
			return &posts, nil // Returned from cache
		}
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
	go redis.Set(redisKey, jsonDatas, 2*time.Hour)

	return &jsonDatas, nil
}

type TedditPostInfo struct {
	Metadata TedditPostMetadata      `json:"metadata"`
	Comments [][]TedditCommmentShape `json:"comments"`
}

// sort must be "best" || "top" || "new" || "controversial" || "old" || "qa"
func GetPostInfo(subteddit, id, sort string) (TedditPostInfo, error) {
	// Check If content already cached:
	redisKey := rediskeys.NewKey(rediskeys.TEDDIT_POST, subteddit+id+sort)
	if postInfo, err := redis.Get[TedditPostInfo](redisKey); err == nil {
		console.Log("Posts Returned from cache", console.Neutral)
		return postInfo, nil // Returned from cache
	}

	Url := fmt.Sprintf("https://teddit.net/r/%s/comments/%s/?sort=%s", url.QueryEscape(subteddit), url.QueryEscape(id), sort)
	if !utils.IsValidURL(Url) {
		return TedditPostInfo{}, errors.New("invalid url")
	}

	htmlPage, err := http.Get(Url)
	if err != nil || htmlPage.StatusCode != 200 {
		return TedditPostInfo{}, errors.New("request failed")
	}
	defer htmlPage.Body.Close()

	doc, err := goquery.NewDocumentFromReader(htmlPage.Body)
	if err != nil {
		return TedditPostInfo{}, err
	}

	PostMetadata := GetPostMetadata(doc)
	PostComments := GetPostComments(doc)
	PostInfos := TedditPostInfo{PostMetadata, PostComments}

	// Caching
	go redis.Set(redisKey, PostInfos, time.Hour)

	return PostInfos, nil
}

type TedditPostMetadata struct {
	Author     string `json:"post_author"`
	Title      string `json:"post_title"`
	Created    int64  `json:"post_created"`
	Ups        string `json:"post_ups"`
	NbComments int    `json:"post_nb_comments"`
	Body       string `json:"body_html"`
}

func GetPostMetadata(doc *goquery.Document) TedditPostMetadata {
	PostAuthor := doc.Find("#post > div.info .title .submitted > a").Text()
	PostTitle := doc.Find("#post > div.info .title > a").Text()
	PostUps := doc.Find("#post > div.info .score > span").Text()

	var PostCreated int64
	if creationISO, exist := doc.Find("#post > div.info .title .submitted > span").Attr("title"); exist {
		layout := "Mon, 02 Jan 2006 15:04:05 GMT"
		if t, err := time.Parse(layout, creationISO); err == nil {
			PostCreated = t.Unix()
		}
	}

	var PostNbComments int
	if PostNbCo := doc.Find("#post > div.comments-info > p").Text(); !utils.IsEmptyString(PostNbCo) {
		if parsed := strings.Split(PostNbCo, " "); len(parsed) == 3 {
			if noCo, err := strconv.Atoi(parsed[1]); err == nil {
				PostNbComments = noCo
			}
		}
	}

	var Body string
	for _, bodyType := range []string{"div.image", "div.video", "div.usertext-body"} {
		if rawbody, err := doc.Find(fmt.Sprintf("#post > %s", bodyType)).Html(); err == nil && !utils.IsEmptyString(rawbody) {
			Body = rawbody
			break
		}
	}

	return TedditPostMetadata{PostAuthor, PostTitle, PostCreated, PostUps, PostNbComments, Body}
}

type TedditCommmentShape struct {
	Id       int `json:"id"`       // Temporary, to rebuild tree in frontend
	ParentId int `json:"parentId"` // Temporary, to rebuild tree in frontend

	Created     int64  `json:"created"`
	Ups         string `json:"ups"`
	Body_html   string `json:"body_html"`
	Link_author string `json:"link_author"`
}

var NodesID = []int{}

func GetPostComments(doc *goquery.Document) [][]TedditCommmentShape {
	var comments [][]TedditCommmentShape
	selection := doc.Find("#post > div.comments > .comment")

	var wg sync.WaitGroup
	wg.Add(selection.Length())
	selection.Each(func(i int, s *goquery.Selection) {
		idIdx := len(NodesID)
		NodesID = append(NodesID, 0)

		go func() {
			defer wg.Done()

			result := RecursiveSearch(s, idIdx, 0)
			sort.Slice(result, func(p, q int) bool {
				return result[p].ParentId < result[q].ParentId // sorting by parentId
			})

			comments = append(comments, result)
		}()
	})
	wg.Wait()

	NodesID = []int{}
	return comments
}

func RecursiveSearch(elem *goquery.Selection, idIdx int, parentId int) []TedditCommmentShape {
	NodesID[idIdx]++

	commentId, _ := elem.Attr("id")
	detailsElem := fmt.Sprintf("#%s > details", commentId)

	CoAuthor := elem.Find(detailsElem + " > .meta .author").First().Text()
	CoUps := strings.TrimSuffix(elem.Find(detailsElem+" > .meta .ups").First().Text(), " points")
	BodyHtml, _ := elem.Find(detailsElem + " > .body").First().Html()

	var CoCreated int64
	if creationISO, exist := elem.Find(detailsElem + " > .meta .created").First().Attr("title"); exist {
		layout := "Mon, 02 Jan 2006 15:04:05 GMT"
		if t, err := time.Parse(layout, creationISO); err == nil {
			CoCreated = t.Unix()
		}
	}
	comment := TedditCommmentShape{NodesID[idIdx], parentId, CoCreated, CoUps, BodyHtml, CoAuthor}

	children := elem.Find(detailsElem + " > .comment")
	if children.Length() <= 0 {
		return []TedditCommmentShape{comment}
	}

	var resultChild = []TedditCommmentShape{}
	children.Each(func(i int, s *goquery.Selection) {
		resultChild = append(resultChild, RecursiveSearch(s, idIdx, NodesID[idIdx])...)
	})

	return append(resultChild, comment)
}
