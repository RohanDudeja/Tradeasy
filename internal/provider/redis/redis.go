package redis

import (
	"crypto/rand"
	"github.com/go-redis/redis"
	"math/big"
	"strconv"
	"time"
)

var redisClient = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "",
	DB:       0,
})

func SetValue(key string, value string, expiry time.Duration) error {
	err := redisClient.Set(key, value, expiry).Err()
	if err != nil {
		return err
	}
	return nil
}
func GetRandNum() (string, error) {
	nBig, e := rand.Int(rand.Reader, big.NewInt(8999))
	if e != nil {
		return "", e
	}
	return strconv.FormatInt(nBig.Int64()+1000, 10), nil
}

func GetValue(key string) (string, error) {
	value, argh := redisClient.Get(key).Result()
	if argh != nil {
		return "", argh
	}
	return value, nil
}
