package utils

import (
	"fmt"

	"github.com/redis/go-redis/v9"
)

var client *redis.Client

func SetClient(redisClient *redis.Client) {
	fmt.Println("Setting Redis client")

	client = redisClient
	fmt.Println("Redis client has been initialized")
}

func GetRedisClient() *redis.Client {
	if client == nil {
		fmt.Println("WARNING: Redis client is nil")
	}
	return client
}
