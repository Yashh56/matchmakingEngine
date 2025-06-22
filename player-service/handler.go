package playerservice

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Yashh56/matchmakingEngine/utils"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

var player Player
var ct = context.Background()

func Join_queue(ctx *gin.Context) {

	if err := ctx.ShouldBindJSON(&player); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var redisClient = utils.GetRedisClient()
	if redisClient == nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Redis not initialized"})
		return
	}

	playerBytes, err := json.Marshal(player)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to serialize player"})
		return
	}
	playerString := string(playerBytes)

	err = redisClient.ZAdd(ctx, "queue:solo:asia", redis.Z{
		Score:  float64(player.JoinedAt),
		Member: playerString,
	}).Err()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Redis Error"})
		return
	}

	val, err := redisClient.Get(ct, player.Player_id).Result()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve player from Redis"})
		return
	}

	var retrievedPlayer Player
	err = json.Unmarshal([]byte(val), &retrievedPlayer)
	if err != nil {
		panic(err)
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status": "Waiting in Queue",
		"player": retrievedPlayer,
	})

}
