package redis

import (
	"Tradeasy/config"
	"github.com/go-redis/redis"
	"log"
	"strconv"
	"time"
)

var redisClient *redis.Client

func CreateClient() {
	redisConfig := config.GetConfig().Redis
	redisClient = redis.NewClient(&redis.Options{
		Addr:     redisConfig.Host + ":" + strconv.Itoa(redisConfig.Port),
		Password: redisConfig.Password,
		DB:       redisConfig.DB,
	})
	pong, err := TestClient()
	if err != nil {
		log.Fatalf("Redis: failed to create client: %v\n", err)
	}
	log.Println(pong)
}
func TestClient() (pong string, err error) {
	pong, err = redisClient.Ping().Result()
	return pong, err
}

func SetValue(key string, value string, expiry time.Duration) error {
	err := redisClient.Set(key, value, expiry).Err()
	if err != nil {
		return err
	}
	return nil
}

func GetValue(key string) (string, error) {
	value, err := redisClient.Get(key).Result()
	if err != nil {
		return "", err
	}
	return value, nil
}
