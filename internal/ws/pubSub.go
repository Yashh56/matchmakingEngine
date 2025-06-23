package ws

import (
	"context"
	"encoding/json"
	"log"

	"github.com/Yashh56/matchmakingEngine/pkg/clients"
	"github.com/redis/go-redis/v9"
)

func ListenForMatches(ctx context.Context, rdb *redis.Client, manager *clients.Manager) {
	sub := rdb.Subscribe(ctx, "matchmaking:events")
	ch := sub.Channel()

	for msg := range ch {
		var payload map[string]interface{}
		if err := json.Unmarshal([]byte(msg.Payload), &payload); err != nil {
			log.Println("Invalid match JSON:", err)
			continue
		}

		players, ok := payload["players"].([]interface{})
		if !ok {
			log.Println("Match payload missing players")
			continue
		}

		for _, p := range players {
			playerObj, ok := p.(map[string]interface{})
			if !ok {
				continue
			}
			playerID := playerObj["player_id"].(string)

			if conn := manager.Get(playerID); conn != nil {
				conn.WriteJSON(map[string]interface{}{
					"type":    "match_found",
					"match":   payload,
					"message": "🎯 You’ve been matched!",
				})
				log.Printf("📨 Match sent to player %s\n", playerID)
			}
		}
	}
}
