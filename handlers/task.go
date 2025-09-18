package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Developer-s-Foundry/DF.2.0-slack-notification/repository/postgres"
	red "github.com/Developer-s-Foundry/DF.2.0-slack-notification/repository/redis"
	"github.com/Developer-s-Foundry/DF.2.0-slack-notification/utils"
	"github.com/julienschmidt/httprouter"
	"github.com/slack-go/slack"
)

type TaskHandler struct {
	DB    *postgres.PostgresConn
	R     *red.RedisConn
	Slack *slack.Client
}

func (t *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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
		ExpiresAt:   task.ExpiresAt,
		AssignedTo:  task.AssignedTo,
		Status:      task.Status,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// var wg = new(sync.WaitGroup)
	go func() {
		publisher(utils.ADD_TASK_TO_DB, 1, tsk, t.R)
	}()

	key := fmt.Sprintf("task:%s", tsk.ID)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	payload, err := json.Marshal(&task)
	if err != nil {
		log.Printf("failed to marshal json data: %v\n", err)
		message := map[string]string{"message": http.StatusText(http.StatusInternalServerError)}
		utils.WriteToJson(w, message, http.StatusInternalServerError)
		return
	}

	if err = t.R.Set(ctx, key, payload); err != nil {
		handleCacheError(w, err, "primary set operation failed")
		return
	}

	if err := t.R.Z(ctx, tsk.ID, task.ExpiresAt.Unix()); err != nil {
		log.Printf("failed to add data to cache: %v\n", err)
		return
	}

	log.Println("data added to cache successfully")

	// send message to slack
	go func() {
		message := fmt.Sprintf(
			":memo: *New Task Created!*\n\n*Title:* %s\n*Assigned To:* %s\n*Status:* %s\n*Description*: %s\n*Due:* %s",
			task.Name,
			task.AssignedTo,
			task.Status,
			task.Description,
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

func (t *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	taskID := ps.ByName("id")

	var updates struct {
		Name        *string    `json:"name"`
		Description *string    `json:"description"`
		ExpiresAt   *time.Time `json:"expires_at"`
		AssignedTo  *string    `json:"assigned_to"`
		Status      *string    `json:"status"`
	}

	if err := utils.ReadDataFromJson(r, &updates); err != nil {
		log.Printf("unable to decode json data: %v", err)
		utils.WriteToJson(w, err.Error(), http.StatusBadRequest)
		return
	}

	existingTask, err := t.DB.GetTaskByID(r.Context(), taskID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.WriteToJson(w, "Task not found", http.StatusNotFound)
			return
		}
		log.Printf("unable to get task from db: %v", err)
		utils.WriteToJson(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if updates.Name != nil {
		existingTask.Name = *updates.Name
	}
	if updates.Description != nil {
		existingTask.Description = *updates.Description
	}
	if updates.ExpiresAt != nil {
		existingTask.ExpiresAt = *updates.ExpiresAt
	}
	if updates.AssignedTo != nil {
		existingTask.AssignedTo = *updates.AssignedTo
	}
	if updates.Status != nil {
		existingTask.Status = *updates.Status
	}

	existingTask.UpdatedAt = time.Now().UTC()

	go func() {
		publisher(utils.UPDATE_TASK_IN_DB, 1, existingTask, t.R)
	}()

	key := fmt.Sprintf("task:%s", existingTask.ID)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	payload, err := json.Marshal(&updates)
	if err != nil {
		log.Printf("failed to marshal json data: %v\n", err)
		message := map[string]string{"message": http.StatusText(http.StatusInternalServerError)}
		utils.WriteToJson(w, message, http.StatusInternalServerError)
		return
	}

	if err := t.R.Set(ctx, key, payload); err != nil {
		handleCacheError(w, err, "primary set operation failed")
		return
	}

	if err := t.R.Z(ctx, existingTask.ID, updates.ExpiresAt.Unix()); err != nil {
		log.Printf("failed to add data to cache: %v\n", err)
		return
	}

	log.Println("data added to cache successfully")

	go func() {
		message := fmt.Sprintf(
			":memo: *Task Updated!*\n\n*Title:* %s\n*Assigned To:* %s\n*Status:* %s\n*Description*: %s\n*Due:* %s",
			existingTask.Name,
			existingTask.AssignedTo,
			existingTask.Status,
			existingTask.Description,
			existingTask.ExpiresAt.Format("Jan 02, 2006 15:04 MST"),
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
		Data:       existingTask,
		StatusCode: http.StatusOK,
		Message:    "task updated successfully",
	}
	utils.WriteToJson(w, response, http.StatusOK)
}

func handleCacheError(w http.ResponseWriter, err error, message string) {
	log.Printf("failed to add data to cache: %v: %s\n", err, message)
	response := map[string]string{"message": http.StatusText(http.StatusInternalServerError)}
	utils.WriteToJson(w, response, http.StatusInternalServerError)
}
