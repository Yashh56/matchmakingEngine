package playerservice

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Yashh56/matchmakingEngine/utils"
	"github.com/gin-gonic/gin"
)

type Player struct {
	Player_id string `json:"player_id"`
	MMR       int    `json:"mmr"`
	Region    string `json:"region"`
	Ping      int    `json:"ping"`
	GameMode  string `json:"game_mode"`
	JoinedAt  int64  `json:"joined_at"`
}

var player Player
var ct = context.Background()

func Join_queue(ctx *gin.Context) {

	if err := ctx.ShouldBindJSON(&player); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var redis = utils.GetRedisClient()
	if redis == nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Redis not initialized"})
		return
	}

	playerBytes, err := json.Marshal(player)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to serialize player"})
		return
	}
	playerString := string(playerBytes)

	err = redis.Set(ct, player.Player_id, playerString, 0).Err()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Redis Error"})
		return
	}

	val, err := redis.Get(ct, player.Player_id).Result()
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
