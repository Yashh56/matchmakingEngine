// cmd/main.go
package main

import (
	"flag"
	"time"

	"github.com/Yashh56/matchmakingEngine/internal/clientSim"
	"github.com/Yashh56/matchmakingEngine/utils"
	"github.com/redis/go-redis/v9"
)

func main() {
	client := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})
	utils.SetClient(client)

	player_id := flag.String("player_id", "", "Player's ID")
	mmr := flag.Int("mmr", 0, "Player's MMR")
	region := flag.String("region", "", "Player's region")
	ping := flag.Int("ping", 0, "Player's ping")

	flag.Parse()

	game_mode := "solo"
	joined_at := time.Now().Unix()

	if *player_id == "" || *region == "" || game_mode == "" || *mmr == 0 || *ping == 0 {
		panic("Missing required flags: player_id, mmr, region, ping, or mode")
	}

	clientSim.Join_Queue(*player_id, *mmr, *region, *ping, game_mode, int(joined_at))

	var redisClient = utils.GetRedisClient()

	if redisClient == nil {
		panic("Something went wrong")
	}
}
