package services

import (
	"errors"
	"net/http"
	"teniditter-server/cmd/global/utils"

	"github.com/PuerkitoBio/goquery"
)

func GetHTMLDocument(URL string) (*goquery.Document, error) {
	if !utils.IsValidURL(URL) {
		return nil, errors.New("invalid URL")
	}

	htmlPage, err := http.Get(URL)
	if err != nil || htmlPage.StatusCode != 200 {
		return nil, errors.New("invalid page")
	}
	defer htmlPage.Body.Close()

	return goquery.NewDocumentFromReader(htmlPage.Body)
}
