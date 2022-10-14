package rediskeys

import (
	"strings"
	"teniditter-server/cmd/global/utils"
)

type RedisKeys string

const (
	TEDDIT_HOME    = RedisKeys("TEDDIT_HOME")
	SUBREDDIT_INFO = RedisKeys("TEDDIT_HOME")
)

func NewKey(base RedisKeys, extendKeyword string) RedisKeys {
	return base + "_" + RedisKeys(strings.ToUpper(utils.FormatToSafeString(extendKeyword)))
}
