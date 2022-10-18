package rediskeys

import (
	"strings"
	"teniditter-server/cmd/global/utils"
)

type RedisKeys string

const (
	TEDDIT_HOME = RedisKeys("TEDDIT_HOME")
	SUBREDDIT   = RedisKeys("TEDDIT_SUBREDDIT")
	USER        = RedisKeys("TEDDIT_USER")
)

func NewKey(base RedisKeys, extendKeyword string) RedisKeys {
	return base + "_" + RedisKeys(strings.ToUpper(utils.FormatToSafeString(extendKeyword)))
}
