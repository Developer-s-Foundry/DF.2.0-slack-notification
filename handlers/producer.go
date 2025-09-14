package handlers

import (
	"encoding/json"
	"log"

	red "github.com/Developer-s-Foundry/DF.2.0-slack-notification/repository/redis"
)

func publisher(topic string, producers int, data interface{}, r *red.RedisConn) {
	for i := 0; i < producers; i++ {
		dByte, err := json.Marshal(data)
		if err != nil {
			log.Printf("unexpected error: %v", err)
			return
		}
		if err := r.Enqueue(topic, dByte); err != nil {
			log.Printf("failed to add data to queue: %v\n", err)
			continue
		}
		log.Printf("pushed task to queue")
	}
}
