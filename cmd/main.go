package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	app "github.com/dimaglushkov/deck-manager/internal"
	"github.com/dimaglushkov/deck-manager/internal/repositories/redis"
)

func run() error {
	var redisHost, redisPort string
	if redisHost = os.Getenv("REDIS_HOST"); redisHost == "" {
		return fmt.Errorf("env variable REDIS_HOST is empty")
	}

	if redisPort = os.Getenv("REDIS_PORT"); redisPort == "" {
		return fmt.Errorf("env variable REDIS_PORT is empty")
	}

	repo, err := redis.New(redisHost, redisPort)
	if err != nil {
		return err
	}
	handler := app.NewHandler(repo)
	return http.ListenAndServe(":"+os.Getenv("APP_PORT"), handler)
}

func main() {
	if err := run(); err != nil {
		log.Fatalf("error while running the app: %s", err.Error())
	}
}
