package websocket

import (
	"context"
	"encoding/json" 
	"log"

	"github.com/centrifugal/gocent/v3"
)

var client = gocent.New(gocent.Config{
	Addr: "http://localhost:8000/api",
	Key:  "PVcA8qZINuBVBlfsZV26GuLlF8hVsy89Lh7ij2aj5-qOg88YpRfdUsdBiQW0HB7ptPTcDbXmcx6LIULMkSFW1w",
})

func PublishToCentrifugo(channel string, data map[string]interface{}) {
	// 1. Convert the map to JSON bytes manually
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Println("Error marshaling data:", err)
		return
	}

	// 2. Pass the byte slice (jsonData) to Publish
	_, err = client.Publish(context.Background(), channel, jsonData)
	
	if err != nil {
		log.Println("Error publishing:", err)
		return
	}
	
	log.Println("Publish succeeded to channel:", channel)
}