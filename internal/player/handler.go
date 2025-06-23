package player

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Yashh56/matchmakingEngine/utils"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

var ct = context.Background()

func Join_queue(ctx *gin.Context) {
	var player Player

	// Bind JSON input
	if err := ctx.ShouldBindJSON(&player); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get Redis client
	redisClient := utils.GetRedisClient()
	if redisClient == nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Redis not initialized"})
		return
	}

	// Serialize player struct to JSON string
	playerBytes, err := json.Marshal(player)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to serialize player"})
		return
	}
	playerString := string(playerBytes)

	// 1. Save player in Redis as key-value pair for future retrieval
	err = redisClient.Set(ct, player.Player_id, playerString, 0).Err()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store player in Redis"})
		return
	}

	// 2. Add player to matchmaking queue (sorted set)
	err = redisClient.ZAdd(ct, "queue:solo:asia", redis.Z{
		Score:  float64(player.JoinedAt),
		Member: playerString, // or just player.Player_id if you want to optimize storage
	}).Err()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add to queue"})
		return
	}

	// 3. Retrieve stored player (optional check)
	val, err := redisClient.Get(ct, player.Player_id).Result()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve player from Redis"})
		return
	}

	// Deserialize JSON back to struct
	var retrievedPlayer Player
	err = json.Unmarshal([]byte(val), &retrievedPlayer)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode player from Redis"})
		return
	}

	// Success response
	ctx.JSON(http.StatusOK, gin.H{
		"status": "Waiting in Queue",
		"player": retrievedPlayer,
	})
}
