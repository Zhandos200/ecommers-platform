package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func NewRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "redis:6379", // Docker hostname
		Password: "",           // No password
		DB:       0,
	})
}

func GetCache(client *redis.Client, key string) (string, error) {
	return client.Get(ctx, key).Result()
}

func SetCache(client *redis.Client, key string, value string, ttl time.Duration) error {
	return client.Set(ctx, key, value, ttl).Err()
}

func DeleteCache(client *redis.Client, key string) error {
	return client.Del(ctx, key).Err()
}
