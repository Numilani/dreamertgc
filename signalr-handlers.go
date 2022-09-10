package main

import (
	"encoding/json"
	"fmt"
	"github.com/philippseith/signalr"
	"hello_go/models"
)

type receiver struct {
	signalr.Hub
	updateChannel chan ServerDataChunk
}

type ServerDataChunk struct {
	callerName string
	data       any
}

func (c *receiver) ReceivePlayerStats(msg string) {
	var stats models.CharacterStatusData
	err := json.Unmarshal([]byte(msg), &stats)
	if err != nil {
		fmt.Println(err)
	}
	c.updateChannel <- ServerDataChunk{
		callerName: "ReceivePlayerStats",
		data:       stats,
	}
}
