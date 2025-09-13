package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

func dispatcher(topic string, workerId int, r *redis.Client, ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		task, err := r.BLPop(ctx, 2*time.Second, topic).Result()
		if err != nil {
			if errors.Is(err, redis.Nil) {
				continue
			}
			log.Printf("worker %d: error reading task: %v", workerId, err)
			time.Sleep(500 * time.Millisecond)
		}
		if len(task) > 1 {
			topic, payload := task[0], task[1]
			log.Printf("worker %d processing task: %s", workerId, payload)
			fmt.Println(payload)
			switch topic {
			// handle each task topic e.g adding to DB or reading to slack get handled from here p;
			case "add-task-to-db":

			}
		}
	}
}

func consumer(topic string, workers int, r *redis.Client) {
	var wg = new(sync.WaitGroup)
	for i := 1; i <= 10; i++ {
		go func(workerId int) {
			wg.Add(1)
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			dispatcher(topic, i, r, ctx, wg)
			cancel()
		}(i)
	}
}
