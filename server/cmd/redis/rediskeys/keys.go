package rediskeys

import (
	"strings"
	"teniditter-server/cmd/global/utils"
)

type RedisKeys string

const (
	TEDDIT_HOME = RedisKeys("TEDDIT_HOME")
	SUBREDDIT   = RedisKeys("TEDDIT_SUBREDDIT")
	TEDDIT_USER = RedisKeys("TEDDIT_USER")
	TEDDIT_POST = RedisKeys("TEDDIT_POST")
)

func NewKey(base RedisKeys, extendKeyword string) RedisKeys {
	return base + "_" + RedisKeys(strings.ToUpper(utils.FormatToSafeString(extendKeyword)))
}