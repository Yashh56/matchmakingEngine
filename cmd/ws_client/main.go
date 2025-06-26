package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

// ANSI color codes for pretty printing
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorGray   = "\033[37m"
	ColorBold   = "\033[1m"
)

func main() {
	playerId := flag.String("player_id", "", "Player ID for WebSocket connection")
	flag.Parse()

	if *playerId == "" {
		log.Fatal("❌ player_id is required")
	}

	printHeader()

	// Connect to WebSocket server
	url := fmt.Sprintf("ws://localhost:8080/ws?player_id=%s", *playerId)
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)

	if err != nil {
		log.Fatalf("❌ Failed to connect to WebSocket: %v\n", err)
	}
	defer conn.Close()

	fmt.Printf("%s✅ Connected to WebSocket as player %s%s\n\n", ColorGreen, *playerId, ColorReset)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			var msg map[string]interface{}
			err := conn.ReadJSON(&msg)
			if err != nil {
				fmt.Printf("%s❌ Read error: %v%s\n", ColorRed, err, ColorReset)
				return
			}

			displayMessage(msg)
		}
	}()

	<-interrupt
	fmt.Printf("\n%s👋 Exiting WebSocket client.%s\n", ColorYellow, ColorReset)
}

func printHeader() {
	fmt.Printf("%s%s", ColorCyan, ColorBold)
	fmt.Println("╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║                    🎮 Game WebSocket Client                  ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")
	fmt.Printf("%s", ColorReset)
	fmt.Println()
}

func displayMessage(msg map[string]interface{}) {
	timestamp := time.Now().Format("15:04:05")

	// Print message header with timestamp
	fmt.Printf("%s%s📨 [%s] New Message Received%s\n", ColorBold, ColorBlue, timestamp, ColorReset)
	fmt.Printf("%s%s%s\n", ColorBlue, strings.Repeat("─", 60), ColorReset)

	// Check message type for special formatting
	if msgType, exists := msg["type"]; exists {
		switch msgType {
		case "match_found":
			displayMatchFound(msg)
		default:
			displayGenericMessage(msg)
		}
	} else {
		displayGenericMessage(msg)
	}

	fmt.Printf("%s%s%s\n\n", ColorBlue, strings.Repeat("─", 60), ColorReset)
}

func displayMatchFound(msg map[string]interface{}) {
	fmt.Printf("%s🎯 MATCH FOUND!%s\n\n", ColorGreen+ColorBold, ColorReset)

	// Display match ID
	if id, exists := msg["id"]; exists {
		fmt.Printf("%s📋 Match ID:%s %s\n", ColorCyan, ColorReset, id)
	}

	// Display game server address
	if address, exists := msg["address"]; exists {
		fmt.Printf("%s🌐 Server:%s %s\n", ColorCyan, ColorReset, address)
	}

	// Display region
	if region, exists := msg["region"]; exists {
		fmt.Printf("%s🗺️  Region:%s %s\n", ColorCyan, ColorReset, strings.ToUpper(region.(string)))
	}

	// Display players information
	if players, exists := msg["players"]; exists {
		fmt.Printf("%s👥 Players:%s\n", ColorCyan, ColorReset)
		displayPlayers(players)
	}

	// Display message text if exists
	if message, exists := msg["message"]; exists {
		fmt.Printf("%s💬 Message:%s %s\n", ColorCyan, ColorReset, message)
	}
}

func displayPlayers(players interface{}) {
	if playerSlice, ok := players.([]interface{}); ok {
		for i, player := range playerSlice {
			if playerMap, ok := player.(map[string]interface{}); ok {
				fmt.Printf("%s   Player %d:%s\n", ColorYellow, i+1, ColorReset)

				// Player ID
				if playerID, exists := playerMap["player_id"]; exists {
					fmt.Printf("     🆔 ID: %v\n", playerID)
				}

				// MMR
				if mmr, exists := playerMap["mmr"]; exists {
					fmt.Printf("     ⭐ MMR: %v\n", mmr)
				}

				// Ping
				if ping, exists := playerMap["ping"]; exists {
					pingValue := ping.(float64)
					pingColor := getPingColor(pingValue)
					fmt.Printf("     📶 Ping: %s%.0fms%s\n", pingColor, pingValue, ColorReset)
				}

				// Game mode
				if gameMode, exists := playerMap["game_mode"]; exists {
					fmt.Printf("     🎮 Mode: %v\n", gameMode)
				}

				// Region
				if region, exists := playerMap["region"]; exists {
					fmt.Printf("     🌍 Region: %v\n", region)
				}

				// Joined at (convert timestamp to readable format)
				if joinedAt, exists := playerMap["joined_at"]; exists {
					if timestamp, ok := joinedAt.(float64); ok {
						joinTime := time.Unix(int64(timestamp), 0)
						fmt.Printf("     ⏰ Joined: %s\n", joinTime.Format("15:04:05"))
					}
				}

				if i < len(playerSlice)-1 {
					fmt.Println()
				}
			}
		}
	}
}

func displayGenericMessage(msg map[string]interface{}) {
	for key, value := range msg {
		displayKeyValue(key, value, 0)
	}
}

func displayKeyValue(key string, value interface{}, indent int) {
	indentStr := strings.Repeat("  ", indent)

	switch v := value.(type) {
	case map[string]interface{}:
		fmt.Printf("%s%s%s:%s\n", indentStr, ColorCyan, key, ColorReset)
		for subKey, subValue := range v {
			displayKeyValue(subKey, subValue, indent+1)
		}
	case []interface{}:
		fmt.Printf("%s%s%s:%s\n", indentStr, ColorCyan, key, ColorReset)
		for i, item := range v {
			displayKeyValue(fmt.Sprintf("[%d]", i), item, indent+1)
		}
	case string:
		fmt.Printf("%s%s%s:%s %s\n", indentStr, ColorCyan, key, ColorReset, v)
	case float64:
		// Check if it's a timestamp
		if key == "joined_at" && v > 1000000000 {
			joinTime := time.Unix(int64(v), 0)
			fmt.Printf("%s%s%s:%s %v (%s)\n", indentStr, ColorCyan, key, ColorReset, v, joinTime.Format("15:04:05"))
		} else {
			fmt.Printf("%s%s%s:%s %v\n", indentStr, ColorCyan, key, ColorReset, v)
		}
	default:
		// Try to pretty print JSON for complex types
		if jsonBytes, err := json.MarshalIndent(v, indentStr+"  ", "  "); err == nil {
			fmt.Printf("%s%s%s:%s\n%s%s\n", indentStr, ColorCyan, key, ColorReset, indentStr, string(jsonBytes))
		} else {
			fmt.Printf("%s%s%s:%s %v\n", indentStr, ColorCyan, key, ColorReset, v)
		}
	}
}

func getPingColor(ping float64) string {
	switch {
	case ping < 50:
		return ColorGreen
	case ping < 100:
		return ColorYellow
	case ping < 200:
		return "\033[33m" // Orange
	default:
		return ColorRed
	}
}
