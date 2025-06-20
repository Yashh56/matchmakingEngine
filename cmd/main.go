// cmd/main.go
package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/Yashh56/matchmakingEngine/internal"
)

func main() {
	fmt.Println("Yash is here")
	player_id := flag.String("player_id", "", "Player's ID")
	mmr := flag.Int("mmr", 0, "Player's MMR")
	region := flag.String("region", "", "Player's region")
	ping := flag.Int("ping", 0, "Player's ping")
	game_mode := flag.String("mode", "", "Game mode")

	flag.Parse()

	joined_at := time.Now().Unix()

	if *player_id == "" || *region == "" || *game_mode == "" || *mmr == 0 || *ping == 0 {
		panic("Missing required flags: player_id, mmr, region, ping, or mode")
	}

	internal.Join_Queue(*player_id, *mmr, *region, *ping, *game_mode, int(joined_at))

}
