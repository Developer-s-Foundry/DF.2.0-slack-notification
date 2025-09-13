package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/Developer-s-Foundry/DF.2.0-slack-notification/utils"
	"github.com/redis/go-redis/v9"
)

func dispatcher(topic string, workerId int, r *redis.Client) {
	for {
		task, err := r.BLPop(context.Background(), 2*time.Second, topic).Result()
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
			switch topic {
			// handle each task topic e.g adding to DB or reading to slack get handled from here p;
			case utils.ADD_TASK_TO_DB:
				fmt.Println(topic, payload)
			}
		}
	}
}

func consumer(topic string, workers int, r *redis.Client) {
	var wg = new(sync.WaitGroup)
	for i := 1; i <= workers; i++ {
		wg.Add(1)
		go func(workerId int) {
			defer wg.Done()
			dispatcher(topic, i, r)
		}(i)
	}
	wg.Wait()
}
