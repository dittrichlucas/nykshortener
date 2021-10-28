package store

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

type StorageService struct {
	redisClient *redis.Client
}

var (
	storeService = &StorageService{}
)

const CacheDuration = 6 * time.Hour

func Redis() *StorageService {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	pong, err := redisClient.Ping().Result()
	if err != nil {
		panic(fmt.Sprintf("Error init Redis: %v", err))
	}

	fmt.Printf("\nRedis started succesfully: pong message = {%s}\n", pong)
	storeService.redisClient = redisClient
	return storeService
}

func SaveURLMapping(shortURL, originalURL string) {
	err := storeService.redisClient.Set(shortURL, originalURL, CacheDuration).Err()
	if err != nil {
		panic(fmt.Sprintf("failed saving key url | Error: %v - shortUrl: %s - originalUrl: %s\n", err, shortURL, originalURL))
	}
}

func RetrieveInitialURL(shortURL string) string {
	result, err := storeService.redisClient.Get(shortURL).Result()
	if err != nil {
		panic(fmt.Sprintf("failed RetrieveInitialURL url | Error: %v - shortUrl: %s\n", err, shortURL))
	}
	return result
}
