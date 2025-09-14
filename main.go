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
	host, port := os.Getenv("DB_HOST"), os.Getenv("DB_PORT")
	password, port := os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT")
	db_name, db_ssl := os.Getenv("DB_NAME"), os.Getenv("DB_SSL")
	redConn, password := os.Getenv("RED_CONN_STRING"), os.Getenv("RED_PASSWORD")
	slackToken := os.Getenv("SLACK_BOT_TOKEN")

	post, err := postgres.ConnectPostgres(url, password, port, host, db_name, user, db_ssl)
	if err != nil {
		panic(err)
	}

	reds, err := red.ConnectRedis(redConn, password, 0)
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
	go consumer(utils.ADD_TASK_TO_DB, 3, reds, post)
	go notifyExpiredTasks(time.Second*5, quitCh, reds)

	// handler registries:
	task := handlers.TaskHandler{DB: post, R: reds, Slack: slk}
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v1/task", task.CreateTask)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", portInt),
		ReadTimeout:  time.Minute * 30,
		WriteTimeout: time.Minute * 30,
		Handler:      mux,
	}
	log.Printf("Server is running on %s\n", server.Addr)
	log.Fatal(server.ListenAndServe())

	<-quitCh
}
