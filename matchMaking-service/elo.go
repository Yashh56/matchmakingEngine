package matchmakingservice

import (
	"context"
	"encoding/json"
	"math"
	"time"

	playerservice "github.com/Yashh56/matchmakingEngine/player-service"
	"github.com/Yashh56/matchmakingEngine/utils"
	"github.com/google/uuid"
)

type Match struct {
	Id      string                 `json:"id"`
	Players []playerservice.Player `json:"players"`
	Region  string                 `json:"region"`
}

func CanMatch(p1, p2 playerservice.Player) bool {
	if p1.Region != p2.Region {
		return false
	}
	mmrGap := math.Abs(float64(p1.MMR - p2.MMR))
	maxGap := 100 + float64(time.Now().Unix()-p1.JoinedAt)/30*10

	return mmrGap <= maxGap && p1.Ping <= 100 && p2.Ping <= 100
}

func FormMatch(p1, p2 playerservice.Player) {
	matchId := uuid.New().String()
	ctx := context.Background()

	match := Match{
		Id:      matchId,
		Players: []playerservice.Player{p1, p2},
		Region:  p1.Region,
	}

	jsonData, _ := json.Marshal(match)
	utils.GetRedisClient().Publish(ctx, "matchmaking:events", jsonData)
}
