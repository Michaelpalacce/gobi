package redis

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	return value
}

func getEnvOrDefaultInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}

	return intValue
}

var rdb *redis.Client = redis.NewClient(&redis.Options{
	Addr:     getEnvOrDefault("REDIS_HOST", "localhost:6379"),
	Password: getEnvOrDefault("REDIS_PASSWORD", ""),
	DB:       getEnvOrDefaultInt("REDIS_DB", 0),
})

var ctx = context.Background()

func Get(key string) (string, error) {
	return rdb.Get(ctx, key).Result()
}

// Set will set a key in redis with a value and expiration
// If the key already exists, it will overwrite the value
// If you don't want to set an expiration, pass 0
func Set(key string, value interface{}, expiration time.Duration) error {
	return rdb.Set(ctx, key, value, expiration).Err()
}

func Del(keys ...string) error {
	return rdb.Del(ctx, keys...).Err()
}

func Expire(key string, expiration time.Duration) error {
	return rdb.Expire(ctx, key, expiration).Err()
}

func Exists(key string) (bool, error) {
	result, err := rdb.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}

	return result == 1, nil
}

func Keys(pattern string) ([]string, error) {
	return rdb.Keys(ctx, pattern).Result()
}

func FlushAll() error {
	return rdb.FlushAll(ctx).Err()
}

func Publish(channel string, message interface{}) error {
	return rdb.Publish(ctx, channel, message).Err()
}

func Subscribe(channels ...string) *redis.PubSub {
	return rdb.Subscribe(ctx, channels...)
}

func Lock(key string, expiration time.Duration) (bool, error) {
	return rdb.SetNX(ctx, key, true, expiration).Result()
}

func Renew(key string, expiration time.Duration) error {
	return Expire(key, expiration)
}

func Unlock(key string) error {
	return Del(key)
}
