package store

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

type StorageService struct {
	redisClient *redis.Client
}

var (
	storeService = &StorageService{}
	ctx          = context.Background()
)

const CacheDuration = 6 * time.Hour

func Redis() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	return client
}

func SaveURLMapping(shortURL, originalURL string) {
	err := storeService.redisClient.Set(ctx, shortURL, originalURL, CacheDuration).Err()
	if err != nil {
		panic(fmt.Sprintf("failed saving key url | Error: %v - shortUrl: %s - originalUrl: %s\n", err, shortURL, originalURL))
	}
}

func RetrieveInitialURL(shortURL string) string {
	result, err := storeService.redisClient.Get(ctx, shortURL).Result()
	if err != nil {
		panic(fmt.Sprintf("failed RetrieveInitialURL url | Error: %v - shortUrl: %s\n", err, shortURL))
	}
	return result
}
