package nitter_test

import (
	"net/url"
	"teniditter-server/cmd/services/nitter"
	"testing"
)

func TestSearchTweetsXML(t *testing.T) {
	if res, err := nitter.SearchTweetsXML(url.QueryEscape("#mikayuu")); err != nil || len(res) <= 0 {
		t.Fatal("SearchTweetsXML() failed", err, res)
	}
}
func TestSearchTweetsScrap(t *testing.T) {
	if res, err := nitter.SearchTweetsScrap("#mikayuu", 1); err != nil || len(res) <= 0 {
		t.Error("SearchTweetsXML() failed", err, res)
	}
}
