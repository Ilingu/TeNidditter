package redis

import (
	"errors"
	"teniditter-server/cmd/global/console"
	"time"

	"github.com/go-redis/redis/v9"
)

var redisConn *redis.Client

func openRedis() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	return rdb
}

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

func Get[T any](key string) (T, error) {
	if redisConn == nil {
		return *new(T), errors.New("no redis conn")
	}

	return *new(T), nil
}

func Set(key string, data any) bool {
	return redisConn != nil
}
