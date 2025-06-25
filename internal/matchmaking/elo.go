package matchmaking

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"time"

	"github.com/Yashh56/matchmakingEngine/internal/player"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type Match struct {
	Id      string          `json:"id"`
	Players []player.Player `json:"players"`
	Region  string          `json:"region"`
}

func CanMatch(p1, p2 player.Player) bool {
	if p1.Region != p2.Region {
		fmt.Printf("[Skip] Region mismatch: %s vs %s\n", p1.Region, p2.Region)
		return false
	}

	mmrGap := math.Abs(float64(p1.MMR - p2.MMR))
	waitTime := float64(time.Now().Unix() - p1.JoinedAt)
	maxGap := 100 + (waitTime/30)*10

	if mmrGap > maxGap {
		fmt.Printf("[Skip] MMR gap too high (%f > %f) between %s and %s\n", mmrGap, maxGap, p1.Player_id, p2.Player_id)
		return false
	}

	if p1.Ping > 100 || p2.Ping > 100 {
		fmt.Printf("[Skip] Ping too high: %s (%d ms), %s (%d ms)\n", p1.Player_id, p1.Ping, p2.Player_id, p2.Ping)
		return false
	}

	fmt.Printf("[Match ✅] %s vs %s | MMR Gap: %.0f | Region: %s\n", p1.Player_id, p2.Player_id, mmrGap, p1.Region)
	return true
}

func FormMatch(p1, p2 player.Player, redisClient *redis.Client) {
	matchId := uuid.New().String()
	ctx := context.Background()

	match := Match{
		Id:      matchId,
		Players: []player.Player{p1, p2},
		Region:  p1.Region,
	}

	jsonData, err := json.Marshal(match)
	if err != nil {
		fmt.Printf("[❌ ERROR] Failed to marshal match: %v\n", err)
		return
	}

	err = redisClient.Publish(ctx, "matchmaking:events", jsonData).Err()
	if err != nil {
		fmt.Printf("[❌ ERROR] Failed to publish match: %v\n", err)
		return
	}

	fmt.Printf("[📣 MATCH PUBLISHED] Match ID: %s | Players: [%s, %s] | Region: %s\n",
		matchId, p1.Player_id, p2.Player_id, p1.Region)
}
