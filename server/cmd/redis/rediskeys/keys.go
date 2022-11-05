package rediskeys

import (
	"strings"
	"teniditter-server/cmd/global/utils"
)

type RedisKeys string

const (
	TEDDIT_HOME      = RedisKeys("TEDDIT_HOME")
	SUBREDDIT        = RedisKeys("TEDDIT_SUBREDDIT")
	TEDDIT_USER      = RedisKeys("TEDDIT_USER")
	TEDDIT_POST      = RedisKeys("TEDDIT_POST")
	TEDDIT_USER_FEED = RedisKeys("TEDDIT_USER_FEED")

	NITTER_NEET_COMMENTS = RedisKeys("NITTER_NEET_COMMENTS")
	NITTER_SEARCH_TWEETS = RedisKeys("NITTER_SEARCH_TWEETS")
	NITTER_NITTOS_TWEETS = RedisKeys("NITTER_NITTOS_TWEETS")
	NITTER_SEARCH_NITTOS = RedisKeys("NITTER_SEARCH_NITTOS")
	NITTER_USER_FEED     = RedisKeys("NITTER_USER_FEED")
)

func NewKey(base RedisKeys, extendKeyword string) RedisKeys {
	return base + "_" + RedisKeys(strings.ToUpper(utils.SafeString(extendKeyword)))
}
