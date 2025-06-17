package main

import (
	"context"
	"fmt"

	playerservice "github.com/Yashh56/matchmakingEngine/player-service"
	"github.com/Yashh56/matchmakingEngine/utils"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func main() {
	router := gin.Default()
	gin.SetMode(gin.ReleaseMode)
	client := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})
	_, err := client.Ping(ctx).Result()

	if err != nil {
		panic(err)
	}

	fmt.Println("Redis has been connected")
	utils.SetClient(client)
	router.POST("/join_queue", playerservice.Join_queue)
	router.Run(":8080")

}
