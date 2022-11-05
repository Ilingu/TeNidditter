package nitter

import (
	"errors"
	"net/url"
	"strings"
	"teniditter-server/cmd/services"

	"github.com/PuerkitoBio/goquery"
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

func fetchTweets(URL string, limit int) ([]NeetComment, error) {
	_, tweetsSelectors := queryMoreSelectors(URL, ".timeline-item", "div.timeline > div.show-more > a", limit)
	if tweetsSelectors == nil {
		return nil, errors.New("no tweets found")
	}

	Tweets := make([]NeetComment, tweetsSelectors.Length())
	tweetsSelectors.Each(func(i int, s *goquery.Selection) {
		Tweets[i] = extractNeetDatas(s)
	})
	return Tweets, nil
}