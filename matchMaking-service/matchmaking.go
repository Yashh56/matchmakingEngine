package matchmakingservice

import (
	"context"
	"encoding/json"

	playerservice "github.com/Yashh56/matchmakingEngine/player-service"
	"github.com/Yashh56/matchmakingEngine/utils"
)

func RunMatchmaking(ctx context.Context) {

	var redisClient = utils.GetRedisClient()

	for {
		matchCandidates, _ := redisClient.ZRange(ctx, "queue:solo:asia", 0, -1).Result()

		for i := 0; i < len(matchCandidates); i++ {
			var player playerservice.Player
			err := json.Unmarshal([]byte(matchCandidates[i]), &player)
			if err != nil {
				continue
			}
			for j := i + 1; j < len(matchCandidates); j++ {
				var candidate playerservice.Player
				err = json.Unmarshal([]byte(matchCandidates[j]), &candidate)

				if err != nil {
					continue
				}

				if CanMatch(player, candidate) {
					FormMatch(player, candidate)
					// RemoveFromQueue(player,candidate)
				}
			}
		}

	}
}
