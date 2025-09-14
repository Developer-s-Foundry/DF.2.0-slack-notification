package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Developer-s-Foundry/DF.2.0-slack-notification/repository/postgres"
	red "github.com/Developer-s-Foundry/DF.2.0-slack-notification/repository/redis"
	"github.com/Developer-s-Foundry/DF.2.0-slack-notification/utils"
	"github.com/slack-go/slack"
)

type TaskHandler struct {
	DB    *postgres.PostgresConn
	R     *red.RedisConn
	Slack *slack.Client
}

func (t *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var task struct {
		Name        string    `json:"name"`
		Description string    `json:"description"`
		ExpiresAt   time.Time `json:"expires_at"`
		AssignedTo  string    `json:"assigned_to"`
		Status      string    `json:"status"`
	}

	if err := utils.ReadDataFromJson(r, &task); err != nil {
		log.Printf("unable to decode json data: %v", err)
		utils.WriteToJson(w, err.Error(), http.StatusBadRequest)
		return
	}

	tsk := postgres.Task{
		ID:          utils.Uuid(),
		Name:        task.Name,
		Description: task.Description,
		Expires_at:  task.ExpiresAt,
		AssignedTo:  task.AssignedTo,
		Status:      task.Status,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	// var wg = new(sync.WaitGroup)
	go func() {
		publisher(utils.ADD_TASK_TO_DB, 1, tsk, t.R)
	}()

	key := fmt.Sprintf("task:%s", tsk.ID)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	payload, err := json.Marshal(&task)
	if err != nil {
		log.Printf("failed to marshal json data: %v\n", err)
		message := map[string]string{"message": http.StatusText(http.StatusInternalServerError)}
		utils.WriteToJson(w, message, http.StatusInternalServerError)
		return
	}

	err = t.R.Set(ctx, key, payload)
	if err != nil {
		log.Printf("failed to add data to cache: %v\n", err)
		message := map[string]string{"message": http.StatusText(http.StatusInternalServerError)}
		utils.WriteToJson(w, message, http.StatusInternalServerError)
		return
	}

	log.Println("data added to cache successfully")

	// send message to slack
	go func() {
		message := fmt.Sprintf(
			":memo: *New Task Created!*\n\n*Title:* %s\n*Assigned To:* %s\n*Status:* %s\n*Due:* %s",
			task.Name,
			task.AssignedTo,
			task.Status,
			task.ExpiresAt.Format("Jan 02, 2006 15:04 MST"),
		)

		if err := utils.SendSlackNotification(t.Slack, message); err != nil {
			log.Printf("failed to send message to slack: %v", err)
			return
		}
	}()
	response := struct {
		Data       interface{} `json:"data"`
		StatusCode int         `json:"status_code"`
		Message    string      `json:"message"`
	}{
		Data:       tsk,
		StatusCode: http.StatusCreated,
		Message:    "task created successfully",
	}
	utils.WriteToJson(w, response, http.StatusCreated)
}
