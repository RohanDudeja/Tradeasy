package redis

import (
	"Tradeasy/config"
	"github.com/go-redis/redis"
	"time"
)

var REDIS *redis.Client

func CreateClient() {
	redisConfig := config.GetConfig().Redis
	REDIS = redis.NewClient(&redis.Options{
		Addr:     redisConfig.Host + ":" + string(redisConfig.Port),
		Password: redisConfig.Password,
		DB:       redisConfig.DB,
	})
}
func TestClient() (pong string, err error) {
	pong, err = REDIS.Ping().Result()
	return pong, err
}
func SetValue(key string, value string, expiry time.Duration) error {
	err := REDIS.Set(key, value, expiry).Err()
	if err != nil {
		return err
	}
	return nil
}

func GetValue(key string) (string, error) {
	value, err := REDIS.Get(key).Result()
	if err != nil {
		return "", err
	}
	return value, nil
}
