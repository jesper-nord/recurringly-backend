package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jesper-nord/recurringly-backend/controller"
	"github.com/jesper-nord/recurringly-backend/entity"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
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

	authCtrl := controller.AuthController{Database: db}
	router := mux.NewRouter()
	router.HandleFunc("/api/login", authCtrl.Login).Methods("POST")
	router.HandleFunc("/api/register", authCtrl.Register).Methods("POST")
	router.HandleFunc("/api/refresh", authCtrl.RefreshToken).Methods("POST")

	authRouter := router.PathPrefix("/api/auth").Subrouter()
	authRouter.Use(controller.JwtMiddleware)

	ctrl := controller.TaskController{Database: db}
	authRouter.HandleFunc("/tasks", ctrl.GetTasks).Methods("GET")
	authRouter.HandleFunc("/tasks/{id}", ctrl.GetTask).Methods("GET")
	authRouter.HandleFunc("/tasks", ctrl.CreateTask).Methods("POST")
	authRouter.HandleFunc("/tasks/{id}", ctrl.EditTask).Methods("PUT")
	authRouter.HandleFunc("/tasks/{id}/complete", ctrl.CompleteTask).Methods("POST")
	authRouter.HandleFunc("/tasks/{id}", ctrl.DeleteTask).Methods("DELETE")
	authRouter.HandleFunc("/tasks/{id}/history/{historyId}", ctrl.EditTaskHistory).Methods("PUT")
	authRouter.HandleFunc("/tasks/{id}/history/{historyId}", ctrl.DeleteTaskHistory).Methods("DELETE")

	clientHost := os.Getenv("CLIENT_HOST")
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{clientHost},
		AllowedMethods:   []string{http.MethodHead, http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		Logger:           log.Default(),
	})

	port := getEnv("PORT", "8090")
	log.Fatal(http.ListenAndServe(":"+port, c.Handler(router)))
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
