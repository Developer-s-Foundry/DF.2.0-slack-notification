package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Developer-s-Foundry/DF.2.0-slack-notification/handlers"
	"github.com/Developer-s-Foundry/DF.2.0-slack-notification/repository/postgres"
	red "github.com/Developer-s-Foundry/DF.2.0-slack-notification/repository/redis"
	"github.com/Developer-s-Foundry/DF.2.0-slack-notification/utils"
	"github.com/Developer-s-Foundry/DF.2.0-slack-notification/utils/seed"
	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
	"github.com/slack-go/slack"
)

func main() {
	quitCh := make(chan struct{})
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT is not found in the environment file!")
	}

	portInt, err := strconv.Atoi(portString)
	if err != nil {
		log.Fatal("Invalid port parameter passed")
	}

	// database setup
	url, user := os.Getenv("DB_URL"), os.Getenv("DB_USER")
	host := os.Getenv("DB_HOST")
	password, port := os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT")
	db_name, db_ssl := os.Getenv("DB_NAME"), os.Getenv("DB_SSL")
	redConn, redPassword := os.Getenv("RED_CONN_STRING"), os.Getenv("RED_PASSWORD")
	slackToken := os.Getenv("SLACK_BOT_TOKEN")

	post, err := postgres.ConnectPostgres(url, password, port, host, db_name, user, db_ssl)
	if err != nil {
		panic(err)
	}

	reds, err := red.ConnectRedis(redConn, redPassword, 0)
	if err != nil {
		log.Printf("redis error: %v\n", err)
		return
	}
	// seed data to DB
	if err := seed.SeedTasks(post); err != nil {
		log.Printf("failed to perform data seeding: %v", err)
	}

	slk := slack.New(slackToken)
	// start consumer queue
	go notifyExpiredTasks(time.Second, quitCh, reds)
	go consumer(utils.ADD_TASK_TO_DB, 5, reds, post, slk)
	go consumer(utils.UPDATE_TASK_IN_DB, 5, reds, post, slk)
	go consumer(utils.NOTIFICATION, 5, reds, post, slk)

	// handler registries:
	task := handlers.TaskHandler{DB: post, R: reds, Slack: slk}
	router := httprouter.New()
	router.POST("/api/v1/task", task.CreateTask)
	router.PATCH("/api/v1/task/:id", task.UpdateTask)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", portInt),
		ReadTimeout:  time.Minute * 30,
		WriteTimeout: time.Minute * 30,
		Handler:      router,
	}
	log.Printf("Server is running on %s\n", server.Addr)
	log.Fatal(server.ListenAndServe())

	<-quitCh
}
