package nitter

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"teniditter-server/cmd/global/utils"
	"teniditter-server/cmd/services"
	"teniditter-server/cmd/services/html"

	"github.com/PuerkitoBio/goquery"
	htmlPkg "golang.org/x/net/html"
)

// This function is error-less so check by yourself if the returned value are nil.
//
// It'll query the selectors of a tweet by scrapping the page and hitting the "show more" button, and repeat this process "limit" times. It then return the concatenated version of all the selectors
func queryMoreSelectors(URL, elemsSelector, nextQuerySelector string, limit int) (pageDoc *goquery.Document, selectors *goquery.Selection) {
	var allSelectors *goquery.Selection

	for i := 0; i < limit; i++ {
		doc, err := services.GetHTMLDocument(URL)
		if err != nil {
			return nil, nil
		}

		selector := doc.Find(elemsSelector)
		if allSelectors == nil {
			pageDoc = doc
			allSelectors = selector
		} else {
			allSelectors = allSelectors.AddNodes(selector.Nodes...)
		}

		nextQuery, exist := doc.Find(nextQuerySelector).Attr("href")
		if !exist {
			return pageDoc, allSelectors // no comments left, return comments already fetched
		}

		nextUrl, err := url.Parse(URL)
		if err != nil {
			return pageDoc, allSelectors
		}

		nextUrl.RawQuery = strings.TrimLeft(nextQuery, "?")
		URL = nextUrl.String()
	}

	return pageDoc, allSelectors
}

func fetchTweets(URL string, limit int) ([][]NeetComment, error) {
	_, tweetsSelectors := queryMoreSelectors(URL, ".timeline-item", "div.timeline > .show-more:not(.timeline-item) > a", limit)
	if tweetsSelectors == nil {
		return nil, errors.New("no tweets found")
	}

	var skip int

	Tweets := [][]NeetComment{}
	tweetsSelectors.Each(func(i int, s *goquery.Selection) {
		if skip > 0 {
			skip--
			return
		}

		thread, toExclude := fetchNeetThread(s)

		skip = toExclude
		Tweets = append(Tweets, thread)
	})
	deleteDuplicatesNeets(&Tweets)

	return Tweets, nil
}

func deleteDuplicatesNeets(neets *[][]NeetComment) {
	seen := map[string]bool{}
	for i, thread := range *neets {
		for j, neet := range thread {
			if _, exist := seen[neet.Id]; exist {
				(*neets)[i] = append(thread[:j], thread[j+1:]...)
			} else {
				seen[neet.Id] = true
			}
		}
	}
}

func GetExternalLinksMetatags(externalUrl string) (map[string]string, error) {
	if !utils.IsValidURL(externalUrl) {
		return nil, errors.New("not valid external url")
	}

	resp, err := http.Get(externalUrl)
	if err != nil || resp.StatusCode != 200 {
		return nil, err
	}
	defer resp.Body.Close()

	htmlBuf, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	htmlPage := string(htmlBuf)

	MetatagsDatas := map[string]string{}
	html.FindElements(htmlPage, "title", func(elem *htmlPkg.Node) (stop bool) {
		MetatagsDatas["title"] = elem.FirstChild.Data
		return true
	})

	acceptedMetatags := []string{"description", "og:image"}
	html.FindElements(htmlPage, "meta", func(elem *htmlPkg.Node) (stop bool) {
		attr := map[string]string{}
		for _, arg := range elem.Attr {
			attr[arg.Key] = arg.Val
		}

		for _, metatag := range acceptedMetatags {
			if attr["name"] == metatag || attr["property"] == metatag {
				MetatagsDatas[metatag] = attr["content"]
			}
		}

		return false
	})

	return MetatagsDatas, nil
}
