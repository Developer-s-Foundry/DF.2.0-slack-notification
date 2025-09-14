package red

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisConn struct {
	RConn *redis.Client
}

func ConnectRedis(host, password string, db int) (*RedisConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	client := redis.NewClient(&redis.Options{
		Addr:     host,
		Password: password,
		DB:       db,
	})

	pong, err := client.Ping(ctx).Result()

	if err != nil {
		log.Printf("Error connecting to redis: %v\n", err)
		return nil, err
	}
	log.Println("REDIS: ", pong)
	return &RedisConn{RConn: client}, nil
}
