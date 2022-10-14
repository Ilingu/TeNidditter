package redis

import (
	"context"
	"encoding/json"
	"errors"
	"teniditter-server/cmd/global/console"
	"teniditter-server/cmd/redis/rediskeys"
	"time"

	"github.com/go-redis/redis/v9"
)

var redisConn *redis.Client

func openRedis() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set (not accessible to clear internet)
		DB:       0,  // use default DB
	})
	return rdb
}

// Will attempt 20 conection to redis
func ConnectRedis() error {
	var i int
	for {
		if i >= 20 {
			console.Log("Failed to Connect to Redis", console.Error, true)
			return errors.New("failed to connect to redis")
		}

		if conn := openRedis(); conn != nil {
			redisConn = conn
			console.Log("Redis Connected Successfully!", console.Success, true)
			return nil
		}

		console.Log("Redis Not Ready Yet, Backing off for 1s...", console.Warning)
		time.Sleep(time.Second)
		i++
	}
}

func DisconnectRedis() error {
	if redisConn == nil {
		return errors.New("no redis conn")
	}

	err := redisConn.Close()
	if err != nil {
		console.Log("Error Cannot Disconnect Redis", console.Error)
	} else {
		console.Log("Redis Disconnected successfully", console.Success)
	}
	return err
}

// Get a value in redis via its key
func Get[T any](key rediskeys.RedisKeys) (T, error) {
	if redisConn == nil {
		return *new(T), errors.New("no redis conn")
	}

	var result T
	rawData, err := redisConn.Get(context.Background(), string(key)).Bytes()
	if err != nil {
		return *new(T), err
	}

	err = json.Unmarshal(rawData, &result)
	// json.NewDecoder(bytes.NewReader(rawData)).Decode(&result)
	if err != nil {
		return *new(T), err
	}

	return result, nil
}

// Set a value in redis cache
func Set(key rediskeys.RedisKeys, data any, exp ...time.Duration) bool {
	if redisConn == nil {
		return false
	}

	jsonBlob, err := json.Marshal(data)
	if err != nil {
		return false
	}

	expiration := 12 * time.Hour
	if len(exp) == 1 {
		expiration = exp[0]
	}

	err = redisConn.Set(context.Background(), string(key), jsonBlob, expiration).Err()
	return err == nil
}
