package xml

import (
	XML "encoding/xml"
)

type Rss[T any] struct {
	Channel Channel[T] `xml:"channel"`
}

type Channel[T any] struct {
	Title string `xml:"title"`
	Link  string `xml:"link"`
	Desc  string `xml:"description"`
	Items []T    `xml:"item"`
}

func ParseRSS[T any](xml []byte) (*Rss[T], error) {
	var parsed Rss[T]
	err := XML.Unmarshal(xml, &parsed)
	if err != nil {
		return nil, err
	}

	return &parsed, nil
}
