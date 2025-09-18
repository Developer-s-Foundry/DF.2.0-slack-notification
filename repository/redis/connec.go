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

func ConnectRedis(redisUrl string) (*RedisConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	opt, _ := redis.ParseURL(redisUrl)
	client := redis.NewClient(opt)
	pong, err := client.Ping(ctx).Result()

	if err != nil {
		log.Printf("Error connecting to redis: %v\n", err)
		return nil, err
	}
	log.Println("REDIS: ", pong)
	return &RedisConn{RConn: client}, nil
}
