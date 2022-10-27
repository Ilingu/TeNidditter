package teddit

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
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

type TedditPostInfo struct {
	Metadata TedditPostMetadata      `json:"metadata"`
	Comments [][]TedditCommmentShape `json:"comments"`
}

func GetPostInfo(subteddit, id string, sort ...string) (TedditPostInfo, error) {
	Sort := "best"
	if len(sort) == 1 && !utils.IsEmptyString(sort[0]) {
		Sort = sort[0]
	}

	// Check If content already cached:
	redisKey := rediskeys.NewKey(rediskeys.TEDDIT_HOME, subteddit+id+Sort)
	if postInfo, err := redis.Get[TedditPostInfo](redisKey); err == nil {
		console.Log("Posts Returned from cache", console.Neutral)
		return postInfo, nil // Returned from cache
	}

	Url := fmt.Sprintf("https://teddit.net/r/%s/comments/%s/?sort=%s", url.QueryEscape(subteddit), url.QueryEscape(id), Sort)
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

	return TedditPostMetadata{PostAuthor, PostTitle, PostCreated, PostUps, PostNbComments}
}

type TedditCommmentShape struct {
	Id       int `json:"id"`       // Temporary, to rebuild tree in frontend
	ParentId int `json:"parentId"` // Temporary, to rebuild tree in frontend

	Created     int64  `json:"created"`
	Ups         int    `json:"ups"`
	Body_html   string `json:"body_html"`
	Link_author string `json:"link_author"`
}

func GetPostComments(doc *goquery.Document) [][]TedditCommmentShape {
	var result [][]TedditCommmentShape
	selection := doc.Find("#post > div.comments > .comment")

	// var wg sync.WaitGroup
	// wg.Add(selection.Length())
	selection.Each(func(i int, s *goquery.Selection) {
		// defer wg.Done()
		result = append(result, RecursiveSearch(s, 0))
	})
	// wg.Wait()

	return result
}

var NodeID = -1

func RecursiveSearch(elem *goquery.Selection, parentId int) []TedditCommmentShape {
	NodeID++
	CoAuthor := elem.Find("details > .meta .author").Text()
	CoUps, _ := strconv.Atoi(elem.Find("details > .meta .ups").Text())
	BodyHtml, _ := elem.Find("details > .body").Html()

	var CoCreated int64
	if creationISO, exist := elem.Find("details > .meta .created").Attr("title"); exist {
		layout := "Mon, 02 Jan 2006 15:04:05 GMT"
		if t, err := time.Parse(layout, creationISO); err == nil {
			CoCreated = t.Unix()
		}
	}
	comment := TedditCommmentShape{NodeID, parentId, CoCreated, CoUps, BodyHtml, CoAuthor}

	children := elem.Find("details > .comment")
	if children.Length() <= 0 {
		return []TedditCommmentShape{comment}
	}

	var resultChild = []TedditCommmentShape{}
	children.Each(func(i int, s *goquery.Selection) {
		resultChild = append(resultChild, RecursiveSearch(s, NodeID)...)
	})

	return append(resultChild, comment)
}
