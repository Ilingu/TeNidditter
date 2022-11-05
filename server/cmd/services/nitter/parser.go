package nitter

import (
	"errors"
	"regexp"
	"strings"
	"teniditter-server/cmd/global/utils"
)

var quoteRegex = regexp.MustCompile(`(?mi:https:\/\/nitter\.pussthecat\.org\/([a-zA-Z]+)\/status\/([0-9]+)#m)`)

func parseNeetUrl(quoteUrl string) (nittos, neetId string) {
	if !strings.HasPrefix(quoteUrl, "https://nitter.pussthecat.org") {
		quoteUrl = "https://nitter.pussthecat.org" + quoteUrl
	}
	if !utils.IsValidURL(quoteUrl) {
		return
	}

	submatchs := quoteRegex.FindStringSubmatch(quoteUrl)
	if len(submatchs) != 3 {
		return
	}

	return submatchs[1], submatchs[2]
}

func fetchQuote(quoteUrl string) (*NeetComment, error) {
	nittos, neetId := parseNeetUrl(quoteUrl)
	if utils.IsEmptyString(nittos) || utils.IsEmptyString(neetId) || len(neetId) != 19 {
		return nil, errors.New("invalid quote url")
	}
	return GetNeetContext(nittos, neetId)
}
