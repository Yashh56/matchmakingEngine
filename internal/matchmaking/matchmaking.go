package matchmaking

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	playerService "github.com/Yashh56/matchmakingEngine/internal/player"
	"github.com/redis/go-redis/v9"
)

func RunMatchmaking(ctx context.Context, redisClient redis.Client) {
	regions := []string{"asia", "europe", "na", "sa"}
	modes := []string{"solo"}

	for {
		for _, mode := range modes {
			for _, region := range regions {

				queueKey := fmt.Sprintf("queue:%s:%s", mode, region)

				matchCandidates, err := redisClient.ZRange(ctx, queueKey, 0, -1).Result()
				if err != nil {
					log.Printf("❌ Error reading from queue %s: %v", queueKey, err)
					continue
				}

				if len(matchCandidates) == 0 {
					log.Printf("ℹ️ No candidates in queue %s", queueKey)
					continue
				}

				matchedSet := make(map[string]bool)

				for i := 0; i < len(matchCandidates); i++ {
					var player playerService.Player
					if err := json.Unmarshal([]byte(matchCandidates[i]), &player); err != nil || matchedSet[player.Player_id] {
						continue
					}

					for j := i + 1; j < len(matchCandidates); j++ {
						var candidate playerService.Player
						if err := json.Unmarshal([]byte(matchCandidates[j]), &candidate); err != nil || matchedSet[candidate.Player_id] {
							continue
						}

						if CanMatch(player, candidate) {
							log.Printf("✅ Match found: %s vs %s in %s (%s)", player.Player_id, candidate.Player_id, region, mode)

							matchedSet[player.Player_id] = true
							matchedSet[candidate.Player_id] = true

							FormMatch(player, candidate, &redisClient)
							RemoveFromQueue(ctx, &redisClient, player, candidate)

							break
						}
					}
				}

			}
		}

		time.Sleep(3 * time.Second)
	}
}

func RemoveFromQueue(ctx context.Context, redisClient *redis.Client, p1, p2 playerService.Player) {
	queueKey := fmt.Sprintf("queue:%s:%s", p1.GameMode, p1.Region)

	p1Json, err1 := json.Marshal(p1)
	P2Json, err2 := json.Marshal(p2)

	if err1 != nil || err2 != nil {
		fmt.Println("[❌ ERROR] Failed to marshal players during ZREM")
		return
	}
	removed, err := redisClient.ZRem(ctx, queueKey, string(p1Json), string(P2Json)).Result()

	err = redisClient.Del(ctx, p1.Player_id).Err()
	err = redisClient.Del(ctx, p2.Player_id).Err()

	if err != nil {
		fmt.Printf("[❌ ERROR] Failed to remove players from queue: %v\n", err)
		return
	}

	fmt.Printf("[🗑️ Removed] %d player(s) removed from queue %s\n", removed, queueKey)

}
