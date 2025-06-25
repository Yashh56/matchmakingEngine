package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Yashh56/matchmakingEngine/internal/matchmaking"
	"github.com/Yashh56/matchmakingEngine/utils"
	"github.com/redis/go-redis/v9"
)

func main() {
	ctx := context.Background()
	client := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})
	utils.SetClient(client)
	for {
		matchmaking.RunMatchmaking(ctx, *utils.GetRedisClient())
		fmt.Println("Live ")
	}
	time.Sleep(3 * time.Second)
}
