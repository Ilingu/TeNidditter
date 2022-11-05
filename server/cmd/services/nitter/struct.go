package nitter

type TweetItem struct {
	Title   string `xml:"title"`
	Creator string `xml:"creator"`
	Desc    string `xml:"description"`
	PubDate string `xml:"pubDate"`
	Guid    string `xml:"guid"`
}
