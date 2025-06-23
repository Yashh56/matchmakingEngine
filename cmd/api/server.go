package main

import (
	"context"
	"fmt"

	"github.com/Yashh56/matchmakingEngine/internal/player"
	"github.com/Yashh56/matchmakingEngine/internal/ws"
	"github.com/Yashh56/matchmakingEngine/pkg/clients"
	"github.com/Yashh56/matchmakingEngine/utils"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func main() {
	router := gin.Default()
	gin.SetMode(gin.ReleaseMode)

	// Redis setup
	client := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})
	if _, err := client.Ping(ctx).Result(); err != nil {
		panic(err)
	}
	fmt.Println("✅ Redis connected")
	utils.SetClient(client)

	// WebSocket client manager and pubsub listener
	clientMgr := clients.NewManager()
	go ws.ListenForMatches(ctx, client, clientMgr)

	// Define HTTP routes
	router.POST("/join_queue", player.Join_queue)
	router.GET("/ws", func(c *gin.Context) {
		ws.HandleWebSocket(clientMgr)(c.Writer, c.Request)
	})

	// Start server
	fmt.Println("🚀 Server running on http://localhost:8080")
	router.Run(":8080")
}
