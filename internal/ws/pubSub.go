package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/Yashh56/matchmakingEngine/internal/gameorchestrator"
	"github.com/Yashh56/matchmakingEngine/pkg/clients"
	"github.com/redis/go-redis/v9"
)

func ListenForMatches(ctx context.Context, rdb *redis.Client, manager *clients.Manager) {
	sub := rdb.Subscribe(ctx, "game:allocated")
	ch := sub.Channel()

	for msg := range ch {
		var payload gameorchestrator.GameSession

		if err := json.Unmarshal([]byte(msg.Payload), &payload); err != nil {
			log.Println("❌ Failed to parse session payload:", err)
			continue
		}

		matchId := payload.MatchId
		address := fmt.Sprintf("%s:%d", payload.Address, payload.Port)

		// 🧠 Get full match data
		matchJSON, err := rdb.Get(ctx, "match:"+matchId).Result()
		if err != nil {
			log.Printf("❌ Could not fetch match data for match %s: %v", matchId, err)
			continue
		}

		var matchData map[string]interface{}
		if err := json.Unmarshal([]byte(matchJSON), &matchData); err != nil {
			log.Printf("❌ Failed to unmarshal match data for match %s: %v", matchId, err)
			continue
		}

		// 🧩 Enrich match data with address and message
		matchData["address"] = address
		matchData["message"] = fmt.Sprintf("🎯 You’ve been matched!\n✅ Pod for match %s created at %s", matchId, address)
		matchData["type"] = "match_found"

		// 🎯 Get matched player IDs
		playerIds, err := rdb.SMembers(ctx, "match_players:"+matchId).Result()
		if err != nil {
			log.Printf("❌ Could not get players for match %s: %v", matchId, err)
			continue
		}

		// 📨 Send to matched players only
		for _, playerId := range playerIds {
			conn := manager.Get(playerId)
			if conn == nil {
				log.Printf("⚠️ No active WebSocket connection for player %s", playerId)
				continue
			}

			if err := conn.WriteJSON(matchData); err != nil {
				log.Printf("❌ Failed to send match to player %s: %v", playerId, err)
			} else {
				log.Printf("📨 Sent full match info to player %s", playerId)
			}
		}
	}
}
