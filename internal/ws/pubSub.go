package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/Yashh56/matchmakingEngine/internal/player"
	"github.com/Yashh56/matchmakingEngine/pkg/clients"
	"github.com/redis/go-redis/v9"
)

type Match struct {
	Id      string          `json:"Id"`
	Players []player.Player `json:"Players"`
	Region  string          `json:"Region"`
}

func ListenForMatches(ctx context.Context, rdb *redis.Client, manager *clients.Manager) {
	sub := rdb.Subscribe(ctx, "matchmaking:events")
	defer sub.Close()

	log.Println("[✅ LISTENER STARTED] Subscribed to 'matchmaking:events' channel")

	ch := sub.Channel()

	for {
		select {
		case <-ctx.Done():
			log.Println("[⚠️ LISTENER STOPPED] Context cancelled, stopping listener")
			return

		case msg, ok := <-ch:
			if !ok {
				log.Println("[⚠️ LISTENER STOPPED] Redis channel closed")
				return
			}

			// Parse the match
			var match Match
			if err := json.Unmarshal([]byte(msg.Payload), &match); err != nil {
				log.Printf("[❌ ERROR] Failed to parse match payload: %v", err)
				continue
			}

			log.Printf("[🎯 MATCH RECEIVED] Match ID: %s | Region: %s", match.Id, match.Region)

			// Build message
			matchData := map[string]interface{}{
				"matchId": match.Id,
				"region":  match.Region,
				"players": match.Players,
				"message": fmt.Sprintf("🎯 You’ve been matched in region %s!", match.Region),
				"type":    "match_found",
			}

			// Send match data to players
			for _, p := range match.Players {
				conn := manager.Get(p.Player_id)
				if conn == nil {
					log.Printf("[⚠️ WARNING] No active WebSocket connection for player %s", p.Player_id)
					continue
				}

				if err := conn.WriteJSON(matchData); err != nil {
					log.Printf("[❌ ERROR] Failed to send match to player %s: %v", p.Player_id, err)
				} else {
					log.Printf("[📨 SENT] Match sent to player %s", p.Player_id)
				}
			}

			// ✅ Stop after handling one message (optional)
			// return
		}
	}
}
