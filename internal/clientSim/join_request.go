package clientSim

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type JoinQueueResponse struct {
	PlayerID string `json:"player_id"`
	MMR      int    `json:"mmr"`
	Region   string `json:"region"`
	Ping     int    `json:"ping"`
	JoinedAt int    `json:"joined_at"`
}

type JoinQueueWrapper struct {
	Status string            `json:"status"`
	Data   JoinQueueResponse `json:"data"`
}

func Join_Queue(player_id string, mmr int, region string, ping int, joined_at int) {
	reqBody := JoinQueueResponse{
		PlayerID: player_id,
		MMR:      mmr,
		Region:   region,
		Ping:     ping,
		JoinedAt: joined_at,
	}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		panic(err)
	}
	res, err := http.Post("http://localhost:8080/join_queue", "application/json", bytes.NewBuffer(jsonData))

	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if err != nil {
		panic(err)
	}
	fmt.Println(string(body))
	var wrapped JoinQueueWrapper
	err = json.Unmarshal(body, &wrapped)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", wrapped.Data)

	fmt.Println(res.Status)
}
