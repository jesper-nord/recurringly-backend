package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jesper-nord/recurringly-backend/controller"
	"github.com/jesper-nord/recurringly-backend/entity"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
)

func main() {
	_ = godotenv.Load()

	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbName := getEnv("DB_NAME", "postgres")
	dbUser := os.Getenv("DB_USERNAME")
	dbPass := os.Getenv("DB_PASSWORD")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPass, dbName)
	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		panic("failed to connect to database")
	}

	// migrate schemas
	err = db.AutoMigrate(&entity.Task{})
	if err != nil {
		panic("failed to migrate schema: 'task'")
	}
	err = db.AutoMigrate(&entity.TaskHistory{})
	if err != nil {
		panic("failed to migrate schema: 'taskHistory'")
	}
	err = db.AutoMigrate(&entity.User{})
	if err != nil {
		panic("failed to migrate schema: 'user'")
	}

	defaultRouter := mux.NewRouter()
	authCtrl := controller.AuthController{Database: db}
	defaultRouter.HandleFunc("/api/login", authCtrl.Login).Methods("POST")
	defaultRouter.HandleFunc("/api/register", authCtrl.Register).Methods("POST")
	defaultRouter.HandleFunc("/api/refresh", authCtrl.RefreshToken).Methods("POST")
	defaultRouter.Use(controller.CorsMiddleware)

	authRouter := defaultRouter.PathPrefix("/api/auth").Subrouter()
	authRouter.Use(controller.JwtMiddleware)
	authRouter.Use(controller.CorsMiddleware)

	ctrl := controller.TaskController{Database: db}
	authRouter.HandleFunc("/tasks", ctrl.GetTasks).Methods("GET")
	authRouter.HandleFunc("/tasks/{id}", ctrl.GetTask).Methods("GET")
	authRouter.HandleFunc("/tasks", ctrl.CreateTask).Methods("POST")
	authRouter.HandleFunc("/tasks/{id}", ctrl.EditTask).Methods("PUT")
	authRouter.HandleFunc("/tasks/{id}/complete", ctrl.CompleteTask).Methods("POST")
	authRouter.HandleFunc("/tasks/{id}", ctrl.DeleteTask).Methods("DELETE")
	authRouter.HandleFunc("/tasks/{id}/history/{historyId}", ctrl.EditTaskHistory).Methods("PUT")
	authRouter.HandleFunc("/tasks/{id}/history/{historyId}", ctrl.DeleteTaskHistory).Methods("DELETE")

	port := getEnv("PORT", "8090")

	log.Fatal(http.ListenAndServe(":"+port, defaultRouter))
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
