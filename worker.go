package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/Developer-s-Foundry/DF.2.0-slack-notification/repository/postgres"
	red "github.com/Developer-s-Foundry/DF.2.0-slack-notification/repository/redis"
	"github.com/Developer-s-Foundry/DF.2.0-slack-notification/utils"
	"github.com/redis/go-redis/v9"
)

func dispatcher(topic string, workerId int, r *red.RedisConn, db *postgres.PostgresConn) {
	for {
		task, err := r.RConn.BLPop(context.Background(), 2*time.Second, topic).Result()
		if err != nil {
			if errors.Is(err, redis.Nil) {
				continue
			}
			log.Printf("worker %d: error reading task: %v", workerId, err)
		}
		if len(task) > 1 {
			topic, payload := task[0], task[1]
			log.Printf("worker %d processing task: %s", workerId, topic)
			switch topic {
			// handle each task topic e.g adding to DB or reading to slack get handled from here p;
			case utils.ADD_TASK_TO_DB:
				var task postgres.Task = postgres.Task{}
				err := json.Unmarshal([]byte(payload), &task)
				if err != nil {
					log.Printf("failed to marshal json data: %v\n", err)
					continue
				}

				err = handleAddToDB(task, db, r)
				if err != nil {
					log.Println(err.Error())
					continue
				}
				// write logic to handle data to databas
			}
		}

		time.Sleep(500 * time.Millisecond)
	}
}

func consumer(topic string, workers int, r *red.RedisConn, db *postgres.PostgresConn) {
	var wg = new(sync.WaitGroup)
	for i := 1; i <= workers; i++ {
		wg.Add(1)
		go func(workerId int) {
			defer wg.Done()
			dispatcher(topic, i, r, db)
		}(i)
	}
	wg.Wait()
}

func handleAddToDB(task postgres.Task, db *postgres.PostgresConn, r *red.RedisConn) error {
	err := db.Insert(task)
	if err != nil {
		log.Printf("failed to add data to db: %v", err)
		return err
	}

	log.Println("data added to DB successfully")

	key := fmt.Sprintf("task:%s", task.ID)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	payload, err := json.Marshal(&task)

	if err != nil {
		log.Printf("failed to marshal json data: %v\n", err)
		return err
	}

	err = r.Set(ctx, key, payload)

	if err != nil {
		log.Printf("failed to add data to cache: %v\n", err)
		return err
	}

	log.Println("data added to cache successfully")
	return nil
}
