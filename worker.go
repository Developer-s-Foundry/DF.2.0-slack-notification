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
	"github.com/slack-go/slack"
)

func dispatcher(topic string, workerId int, r *red.RedisConn, db *postgres.PostgresConn, slk *slack.Client) {
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

				err = handleAddToDB(task, db)
				if err != nil {
					log.Println(err.Error())
					continue
				}
			case utils.NOTIFICATION:
				slackNotificationMessage(slk, r, []byte(payload))
			}
		}

		time.Sleep(500 * time.Millisecond)
	}
}

func consumer(topic string, workers int, r *red.RedisConn, db *postgres.PostgresConn, slk *slack.Client) {
	var wg = new(sync.WaitGroup)
	for i := 1; i <= workers; i++ {
		wg.Add(1)
		go func(workerId int) {
			defer wg.Done()
			dispatcher(topic, i, r, db, slk)
		}(i)
	}
	wg.Wait()
}

func handleAddToDB(task postgres.Task, db *postgres.PostgresConn) error {
	err := db.Insert(task)
	if err != nil {
		log.Printf("failed to add data to db: %v", err)
		return err
	}

	log.Println("data added to DB successfully")
	return nil
}

func slackNotificationMessage(slk *slack.Client, r *red.RedisConn, payload []byte) {
	var Payload struct {
		TaskId string `json:"taskId"`
		Event  string `json:"event"`
	}

	if err := json.Unmarshal(payload, &Payload); err != nil {
		log.Printf("failed to marshal json: %v", err)
		return
	}

	key := "task:" + Payload.TaskId

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var task *postgres.Task = &postgres.Task{}
	if err := r.Get(ctx, key, task); err != nil {
		log.Printf("failed to fetch data from cache: %v", err)
		return
	}

	switch Payload.Event {
	case "task-expires":
		var message string

		if task.Status == "completed" {
			message = fmt.Sprintf(
				":tada: *Task Completed!*\n\n*Title:* %s\n*Assigned To:* %s\n*Status:* %s\n*Completed At:* %s\n\nGreat job! :clap:",
				task.Name,
				task.AssignedTo,
				task.Status,
				task.UpdatedAt.Format("Jan 02, 2006 15:04 MST"),
			)
		} else if time.Now().After(task.Expires_at) {
			message = fmt.Sprintf(
				":warning: *Task Expired!*\n\n*Title:* %s\n*Assigned To:* %s\n*Status:* %s\n*Expired At:* %s\n\nPlease review and take action.",
				task.Name,
				task.AssignedTo,
				task.Status,
				task.Expires_at.Format("Jan 02, 2006 15:04 MST"),
			)
		} else {
			message = fmt.Sprintf(
				":memo: *Task Update!*\n\n*Title:* %s\n*Assigned To:* %s\n*Status:* %s\n*Due At:* %s",
				task.Name,
				task.AssignedTo,
				task.Status,
				task.Expires_at.Format("Jan 02, 2006 15:04 MST"),
			)
		}
		if err := utils.SendSlackNotification(slk, message); err != nil {
			log.Printf("failed to notify user on task expiry: %v", err)
			return
		}
		log.Printf("notification sent to user slack on taskId: %s", key)

		if err := r.Del(key); err != nil {
			log.Printf("could not delete cache key: %v", err)
		}
	}
}

func notifyExpiredTasks(interval time.Duration, quitCh chan struct{}, r *red.RedisConn) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-quitCh:
			return
		case <-ticker.C:
			now := time.Now().Unix()

			tasks, err := r.RConn.ZRangeByScore(context.Background(), red.TaskExpirations, &redis.ZRangeBy{
				Min: "-inf",
				Max: fmt.Sprintf("%f", float64(now)),
			}).Result()

			if err != nil {
				log.Println("Redis error:", err)
				time.Sleep(time.Second * 2)
				continue
			}
			for _, taskID := range tasks {
				r.RConn.ZRem(context.Background(), red.TaskExpirations, taskID)
				notify := map[string]string{
					"taskId": taskID,
					"event":  "task-expires",
				}

				data, _ := json.Marshal(notify)
				err := r.Enqueue(utils.NOTIFICATION, data)

				if err != nil {
					log.Printf("failed to add task to queue: %v", err)
					time.Sleep(time.Second * 2) // sleep for 2 seconds
					continue
				}
				log.Printf("added task notification to queue: %s", taskID)
			}
		}
	}
}
