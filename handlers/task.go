package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/Developer-s-Foundry/DF.2.0-slack-notification/repository/postgres"
	red "github.com/Developer-s-Foundry/DF.2.0-slack-notification/repository/redis"
	"github.com/Developer-s-Foundry/DF.2.0-slack-notification/utils"
)

type TaskHandler struct {
	DB *postgres.PostgresConn
	R  *red.RedisConn
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
