package player

type Player struct {
	Player_id string `json:"player_id"`
	MMR       int    `json:"mmr"`
	Region    string `json:"region"`
	Ping      int    `json:"ping"`
	JoinedAt  int64  `json:"joined_at"`
}
