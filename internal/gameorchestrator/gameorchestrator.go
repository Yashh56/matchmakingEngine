package gameorchestrator

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type Match struct {
	Id      string        `json:"id"`
	Players []interface{} `json:"players"`
	Region  string        `json:"region"`
}

type GameSession struct {
	MatchId   string `json:"match_id"`
	Address   string `json:"address"`
	Port      int    `json:"port"`
	SessionId string `json:"session_id"`
}

func Start(ctx context.Context, rdb *redis.Client) {
	sub := rdb.Subscribe(ctx, "matchmaking:events")
	ch := sub.Channel()
	clientSet, err := NewKubeClient()

	if err != nil {
		log.Fatal("Failed to create Kubernetes client:", err)
	}

	for msg := range ch {
		var match Match
		if err := json.Unmarshal([]byte(msg.Payload), &match); err != nil {
			log.Println("❌ Invalid match payload:", err)
			continue
		}

		go func(m Match) {
			address, err := CreateGamePod(ctx, clientSet, m.Id)
			if err != nil {
				log.Println("❌ Pod creation failed:", err)
				return
			}
			session := GameSession{
				MatchId:   m.Id,
				Address:   address,
				Port:      8080,
				SessionId: uuid.NewString(),
			}
			sessionData, _ := json.Marshal(session)
			rdb.Set(ctx, "game_session"+m.Id, sessionData, 0)
			rdb.Publish(ctx, "game:allocated", sessionData)
			log.Printf("✅ Pod for match %s created at %s\n", m.Id, address)

		}(match)

		// go allocateGameServer(ctx, rdb, match)
	}
}

func allocateGameServer(ctx context.Context, rdb *redis.Client, match Match) {
	address := fmt.Sprintf("game-server-%s.example.com", match.Id[:8])
	port := 3000 + (uuid.New().ID() % 1000)

	session := GameSession{
		MatchId:   match.Id,
		Address:   address,
		Port:      int(port),
		SessionId: uuid.NewString(),
	}

	data, _ := json.Marshal(session)

	err := rdb.Set(ctx, "game_session:"+match.Id, data, 0).Err()

	if err != nil {
		log.Println("❌ Failed to store game session:", err)
		return
	}

	log.Printf("🎮 Game server allocated for match %s → %s:%d\n", match.Id, session.Address, session.Port)

	// Optionally notify another channel
	rdb.Publish(ctx, "game:allocated", data)
}
