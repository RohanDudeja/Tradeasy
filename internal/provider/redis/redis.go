package redis

import (
	"github.com/go-redis/redis"
	"time"
)

func createClient() *redis.Client {
	var redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	return redisClient
}

func SetValue(key string, value string, expiry time.Duration) error {
	err := createClient().Set(key, value, expiry).Err()
	if err != nil {
		return err
	}
	return nil
}

func GetValue(key string) (string, error) {
	value, err := createClient().Get(key).Result()
	if err != nil {
		return "", err
	}
	return value, nil
}
