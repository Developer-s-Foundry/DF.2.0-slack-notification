package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

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

	mux := http.Server{
		Addr:         fmt.Sprintf(":%d", portInt),
		ReadTimeout:  time.Minute * 30,
		WriteTimeout: time.Minute * 30,
	}

	log.Fatal(mux.ListenAndServe())
}
