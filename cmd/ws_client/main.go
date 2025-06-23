package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/gorilla/websocket"
)

func main() {
	playerId := flag.String("player_id", "", "Player ID for WebSocket connection")
	flag.Parse()

	if *playerId == "" {
		log.Fatal("❌ player_id is required")
	}

	// Connect to WebSocket server
	url := fmt.Sprintf("ws://localhost:8080/ws?player_id=%s", *playerId)
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)

	if err != nil {
		log.Fatalf("❌ Failed to connect to WebSocket: %v\n", err)
	}
	defer conn.Close()

	fmt.Printf("✅ Connected to WebSocket as player %s\n", *playerId)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			var msg map[string]interface{}
			err := conn.ReadJSON(&msg)
			if err != nil {
				log.Println("❌ Read error:", err)
				return
			}

			fmt.Println("📨 New WebSocket Message:")
			for k, v := range msg {
				fmt.Printf("  %s: %v\n", k, v)
			}
			fmt.Println()
		}
	}()

	<-interrupt
	fmt.Println("\n👋 Exiting WebSocket client.")

}
