package teddit

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"teniditter-server/cmd/global/console"
	"teniditter-server/cmd/global/utils"
	utils_enc "teniditter-server/cmd/global/utils/encryption"
	"teniditter-server/cmd/redis"
	"teniditter-server/cmd/redis/rediskeys"
	"time"
)

func GetUserInfos(username string) (*map[string]any, error) {
	redisKey := rediskeys.NewKey(rediskeys.TEDDIT_USER, utils_enc.Hash(username))

	if posts, err := redis.Get[map[string]any](redisKey); err == nil {
		console.Neutral("Teddit User Info Returned from Cache ⚡")
		return &posts, nil
	}

	Url := fmt.Sprintf("https://teddit.net/u/%s?api&raw_json=1", url.QueryEscape(username))
	if !utils.IsValidURL(Url) {
		return nil, errors.New("invalid URL")
	}

	rawUserInfo, err := http.Get(Url)
	if err != nil || rawUserInfo.StatusCode != 200 {
		return nil, err
	}
	defer rawUserInfo.Body.Close()

	rawBlobUserInfo, err := io.ReadAll(rawUserInfo.Body)
	if err != nil {
		return nil, err
	}

	var userInfo map[string]any
	err = json.Unmarshal(rawBlobUserInfo, &userInfo)
	if err != nil || len(userInfo) <= 0 {
		return nil, err
	}

	// Caching
	go redis.Set(redisKey, userInfo, 4*24*time.Hour) // 4d

	return &userInfo, nil
}
