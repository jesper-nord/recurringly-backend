package main

import (
	"github.com/jesper-nord/recurringly-backend/auth"
	"github.com/jesper-nord/recurringly-backend/postgres"
	"github.com/jesper-nord/recurringly-backend/router"
	"github.com/jesper-nord/recurringly-backend/task"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func main() {
	_ = godotenv.Load()

	db, err := postgres.Initialize()
	if err != nil {
		log.Panicf("failed to connect to database: %v", err)
	}

	taskRepository := task.NewRepository(db)
	err = taskRepository.Migrate()
	if err != nil {
		log.Panicf("failed to migrate schemas: %v", err)
	}

	authRepository := auth.NewRepository(db)
	err = authRepository.Migrate()
	if err != nil {
		log.Panicf("failed to migrate schemas: %v", err)
	}

	routes := router.New(task.NewService(taskRepository), auth.NewService(authRepository))
	port := getEnvWithFallback("PORT", "8090")
	log.Fatal(http.ListenAndServe(":"+port, routes))
}

func getEnvWithFallback(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
