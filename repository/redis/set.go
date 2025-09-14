package red

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var TaskExpirations = "task_expirations"

func (r *RedisConn) Set(ctx context.Context, key string, data interface{}) error {
	err := r.RConn.Set(ctx, key, data, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisConn) Get(ctx context.Context, key string, receiver interface{}) error {
	res, err := r.RConn.Get(ctx, key).Result()

	if err != nil {
		if errors.Is(err, redis.Nil) {
			return fmt.Errorf("key %q not found in cache", key)
		}
		return err
	}
	switch val := receiver.(type) {
	case *string:
		*val = res
		return nil
	default:
		return json.Unmarshal([]byte(res), receiver)
	}
}

func (r *RedisConn) Del(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err := r.RConn.Del(ctx, key).Err()

	if err != nil {
		return fmt.Errorf("failed to delete key %q: %w", key, err)
	}
	return nil
}

func (r *RedisConn) Z(ctx context.Context, taskId string, expiresAt int64) {
	r.RConn.ZAdd(ctx, TaskExpirations, redis.Z{
		Score:  float64(expiresAt),
		Member: taskId,
	})
}
