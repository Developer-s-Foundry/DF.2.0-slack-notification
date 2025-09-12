package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Developer-s-Foundry/DF.2.0-slack-notification/repository/postgres"
	"github.com/Developer-s-Foundry/DF.2.0-slack-notification/utils/seed"
	"github.com/joho/godotenv"
)

func main() {
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

	url, user := os.Getenv("DB_URL"), os.Getenv("DB_USER")
	host, port := os.Getenv("DB_HOST"), os.Getenv("DB_PORT")
	password, port := os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT")
	db_name, db_ssl := os.Getenv("DB_NAME"), os.Getenv("DB_SSL")

	post, err := postgres.ConnectPostgres(url, password, port, host, db_name, user, db_ssl)

	if err != nil {
		panic(err)
	}

	if !seed.Seeded {
		if err := seed.SeedTasks(post); err != nil {
			log.Printf("failed to perform data seeding: %v", err)
		}
		log.Printf("Seeded %d dummy tasks successfully", len(seed.Data()))
	} else {
		log.Printf("Data already seeded")
	}

	mux := http.Server{
		Addr:         fmt.Sprintf(":%d", portInt),
		ReadTimeout:  time.Minute * 30,
		WriteTimeout: time.Minute * 30,
	}

	log.Fatal(mux.ListenAndServe())
}
