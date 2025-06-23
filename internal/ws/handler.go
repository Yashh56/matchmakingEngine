package ws

import (
	"log"
	"net/http"

	"github.com/Yashh56/matchmakingEngine/pkg/clients"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func HandleWebSocket(manager *clients.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		playerID := r.URL.Query().Get("player_id")
		if playerID == "" {
			http.Error(w, "player_id required", http.StatusBadRequest)
			return
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("WebSocket upgrade error:", err)
			return
		}

		manager.Add(playerID, conn)
		log.Printf("🟢 Player connected: %s\n", playerID)
	}
}
