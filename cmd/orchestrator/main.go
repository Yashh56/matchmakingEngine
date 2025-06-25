package main

import (
	"context"
	"log"

	"github.com/Yashh56/matchmakingEngine/internal/gameorchestrator"
	"github.com/Yashh56/matchmakingEngine/utils"
	"github.com/redis/go-redis/v9"
)

func main() {
	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})

	utils.SetClient(rdb)

	log.Println("🚀 Game Orchestrator started")
	gameorchestrator.Start(ctx, rdb)
}
