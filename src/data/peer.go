package data

import (
	"encoding/json"
	"fmt"
)

type Peer struct {
	ID      string `json:"id"`
	Address string `json:"address"`
}

func PeerFromJson(jsonBytes []byte) Peer {
	p := Peer{}
	err := json.Unmarshal(jsonBytes, &p)
	if err != nil {
		fmt.Println("Error in getting Peer from JSON")
	}
	return p
}
