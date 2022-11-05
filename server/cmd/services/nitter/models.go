package nitter

// XML version of NeetBasicComment
type XmlTweetItem struct {
	Title   string `xml:"title"`
	Creator string `xml:"creator"`
	Desc    string `xml:"description"`
	PubDate string `xml:"pubDate"`
	Guid    string `xml:"guid"`
}

// Tweet
type NeetInfo struct {
	MainThread []NeetComment   `json:"main"`
	Reply      [][]NeetComment `json:"reply"`
}

type NeetComment struct {
	NeetBasicComment
	Quote *NeetBasicComment `json:"quote,omitempty"`
}
type NeetBasicComment struct {
	Content      string           `json:"content"`
	Creator      NittosPreview    `json:"creator"`
	CreatedAt    int              `json:"createdAt"`
	Stats        NeetCommentStats `json:"stats"`
	Attachments  *Attachments     `json:"attachment,omitempty"`
	ExternalLink string           `json:"externalLink,omitempty"`
}
type Attachments struct {
	ImagesUrls []string `json:"images,omitempty"`
	VideosUrls []string `json:"videos,omitempty"`
}
type NeetCommentStats struct {
	ReplyCounts  int `json:"reply_counts"`
	RTCounts     int `json:"rt_counts"`
	QuotesCounts int `json:"quotes_counts,omitempty"`
	LikesCounts  int `json:"likes_counts"`
	PlayCounts   int `json:"play_counts,omitempty"`
}

// User
type NittosPreview struct {
	Username    string `json:"username"`
	Description string `json:"description"`
	AvatarUrl   string `json:"avatarUrl"`
}

type Nittos struct {
	Username  string      `json:"username"`
	Bio       string      `json:"bio"`
	AvatarUrl string      `json:"avatarUrl"`
	Location  string      `json:"location"`
	Website   string      `json:"website"`
	JoinDate  string      `json:"joinDate"`
	Stats     NittosStats `json:"stats"`
	BannerUrl string      `json:"bannerUrl"`
}
type NittosStats struct {
	TweetsCounts    int `json:"tweets_counts"`
	FollowingCounts int `json:"following_counts"`
	FollowersCounts int `json:"followers_counts"`
	LikesCounts     int `json:"likes_counts"`
}
