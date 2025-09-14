package handlers

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/Developer-s-Foundry/DF.2.0-slack-notification/repository/postgres"
	red "github.com/Developer-s-Foundry/DF.2.0-slack-notification/repository/redis"
	"github.com/Developer-s-Foundry/DF.2.0-slack-notification/utils"
	"github.com/julienschmidt/httprouter"
)

type TaskHandler struct {
	DB *postgres.PostgresConn
	R  *red.RedisConn
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
		existingTask.Expires_at = *updates.ExpiresAt
	}
	if updates.AssignedTo != nil {
		existingTask.AssignedTo = *updates.AssignedTo
	}
	if updates.Status != nil {
		existingTask.Status = *updates.Status
	}

	existingTask.UpdatedAt = time.Now().UTC()

	if err := t.DB.UpdateTask(r.Context(), *existingTask); err != nil {
		log.Printf("unable to update task in db: %v", err)
		utils.WriteToJson(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

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
