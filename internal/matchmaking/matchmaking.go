package matchmaking

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	playerservice "github.com/Yashh56/matchmakingEngine/player-service"
	"github.com/redis/go-redis/v9"
)

func RunMatchmaking(ctx context.Context, redisClient redis.Client) {

	for {

		matchCandidates, _ := redisClient.ZRange(ctx, "queue:solo:asia", 0, -1).Result()

		matchedSet := make(map[string]bool)

		for i := 0; i < len(matchCandidates); i++ {
			var player playerservice.Player
			err := json.Unmarshal([]byte(matchCandidates[i]), &player)
			if err != nil || matchedSet[player.Player_id] {
				continue
			}

			for j := i + 1; j < len(matchCandidates); j++ {
				var candidate playerservice.Player
				err = json.Unmarshal([]byte(matchCandidates[j]), &candidate)
				if err != nil || matchedSet[candidate.Player_id] {
					continue
				}

				if CanMatch(player, candidate) {
					// Mark as matched
					matchedSet[player.Player_id] = true
					matchedSet[candidate.Player_id] = true

					// Form the match
					FormMatch(player, candidate, &redisClient)

					// Remove both players from Redis
					RemoveFromQueue(ctx, &redisClient, player, candidate)

					break // stop looking for a match for player i
				}
			}
		}

		// Wait before next matchmaking cycle
		time.Sleep(3 * time.Second)
	}
}

func RemoveFromQueue(ctx context.Context, redisClient *redis.Client, p1, p2 playerservice.Player) {
	queueKey := fmt.Sprintf("queue:%s:%s", p1.GameMode, p1.Region)

	p1Json, err1 := json.Marshal(p1)
	P2Json, err2 := json.Marshal(p2)

	if err1 != nil || err2 != nil {
		fmt.Println("[❌ ERROR] Failed to marshal players during ZREM")
		return
	}
	removed, err := redisClient.ZRem(ctx, queueKey, string(p1Json), string(P2Json)).Result()

	if err != nil {
		fmt.Printf("[❌ ERROR] Failed to remove players from queue: %v\n", err)
		return
	}

	fmt.Printf("[🗑️ Removed] %d player(s) removed from queue %s\n", removed, queueKey)

}
