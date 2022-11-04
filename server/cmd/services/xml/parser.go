package xml

import (
	XML "encoding/xml"
)

func ParseNitterSearch(xml []byte) (*Rss[TweetItem], error) {
	var parsed Rss[TweetItem]
	err := XML.Unmarshal(xml, &parsed)
	if err != nil {
		return nil, err
	}

	return &parsed, nil
}
