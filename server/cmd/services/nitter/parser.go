package nitter

import (
	"errors"
	"net/url"
	"regexp"
	"strings"
	"teniditter-server/cmd/global/utils"
	"teniditter-server/cmd/services/html"
	"time"

	htmlpkg "golang.org/x/net/html"
)

var quoteRegex = regexp.MustCompile(`(?mi:https:\/\/nitter\.pussthecat\.org\/([a-zA-Z]+)\/status\/([0-9]+)#m)`)

func parseNeetUrl(neetUrl string) (nittos, neetId string) {
	if !strings.HasPrefix(neetUrl, "https://nitter.pussthecat.org") {
		neetUrl = "https://nitter.pussthecat.org" + neetUrl
	}
	if !utils.IsValidURL(neetUrl) {
		return
	}

	submatchs := quoteRegex.FindStringSubmatch(neetUrl)
	if len(submatchs) != 3 {
		return
	}

	return submatchs[1], submatchs[2]
}

func fetchCtxNeetFromUrl(neetUrl string) (*NeetComment, error) {
	nittos, neetId := parseNeetUrl(neetUrl)
	if utils.IsEmptyString(nittos) || utils.IsEmptyString(neetId) || len(neetId) != 19 {
		return nil, errors.New("invalid neet url")
	}
	return GetNeetContext(nittos, neetId)
}

// This method convert a Xml queried Tweet struct to the standar tweet struct (JSON) that can be send back to client
//
// This function come in to flavour with the [fromXmlDatas] argument:
//   - if set to `true`: the function will only take datas from the provided Xml tweet, note that [stats] and [videos] are not datas present in the xml form of the tweet, thus they are missing in the converted form
//   - if set to `false`: the function will fetch and scrap all the necessary info of this tweet without caring about the xml datas (except for the Guid field containing the url to fetch the tweet's datas) and return a converted form without missing datas, however this method is slower than the previous one
func (neet XmlTweetItem) ToJSON(fromXmlDatas bool) (*NeetComment, error) {
	if fromXmlDatas {
		// Convert string pubDate to unix timestamp
		var createdAt int64
		const layout = "Mon, 02 Jan 2006 15:04:05 GMT"
		if t, err := time.Parse(layout, neet.PubDate); err == nil {
			createdAt = t.Unix()
		}

		// Find Images
		imgSrc := []string{}
		html.FindElements(neet.Desc, "img", func(elem *htmlpkg.Node) (stop bool) {
			for _, attr := range elem.Attr {
				if attr.Key == "src" {
					imgSrc = append(imgSrc, attr.Val)
				}
			}
			return false
		})

		// Find last external links
		var linkCard string
		html.FindElements(neet.Desc, "a", func(elem *htmlpkg.Node) (stop bool) {
			for _, attr := range elem.Attr {
				if attr.Key == "href" {
					if !utils.IsValidURL(attr.Val) {
						continue
					}
					if link, err := url.Parse(attr.Val); err == nil && !link.IsAbs() {
						linkCard = attr.Val
					}
				}
			}
			return false
		})

		// missing: stats and videos
		commentMetadata := NeetBasicComment{neet.Desc, NittosPreview{Username: neet.Creator}, int(createdAt), NeetCommentStats{}, &Attachments{ImagesUrls: imgSrc}, linkCard}
		comment := NeetComment{NeetBasicComment: commentMetadata}

		return &comment, nil
	}

	return fetchCtxNeetFromUrl(neet.Guid) // refetch all
}
