package red

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

func (r *RedisConn) Enqueue(key string, data interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := r.RConn.RPush(ctx, key, data).Err()

	if err != nil {
		if redis.Nil == err {
			return fmt.Errorf("failed to queue error")
		}
		return err
	}

	return nil
}

func (r *RedisConn) Dequeue(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := r.RConn.BLPop(ctx, 2*time.Second, key).Err()

	if err != nil {
		return err
	}
	return nil
}
