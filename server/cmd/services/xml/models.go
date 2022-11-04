package xml

type Rss[T any] struct {
	Channel Channel[T] `xml:"channel"`
}

type Channel[T any] struct {
	Title string `xml:"title"`
	Link  string `xml:"link"`
	Desc  string `xml:"description"`
	Items []T    `xml:"item"`
}

type TweetItem struct {
	Title   string `xml:"title"`
	Creator string `xml:"dc:creator"`
	Desc    string `xml:"description"`
	PubDate string `xml:"pubDate"`
	Guid    string `xml:"guid"`
}
