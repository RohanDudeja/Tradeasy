package redis

import (
	"Tradeasy/config"
	"github.com/go-redis/redis"
	"time"
)

func createClient() *redis.Client {
	redisConfig := config.GetConfig().Redis
	var redisClient = redis.NewClient(&redis.Options{
		Addr:     redisConfig.Host + ":" + string(redisConfig.Port),
		Password: redisConfig.Password,
		DB:       redisConfig.DB,
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