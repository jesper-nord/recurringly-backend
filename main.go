package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"recurringly-backend/controller"
	"recurringly-backend/entity"
)

func main() {
	_ = godotenv.Load()

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=postgres sslmode=disable", os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"))
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

	authRouter := defaultRouter.PathPrefix("/api/auth").Subrouter()
	authRouter.Use(controller.JwtMiddleware)
	authRouter.HandleFunc("/refresh", authCtrl.RefreshToken).Methods("POST")

	ctrl := controller.TaskController{Database: db}
	authRouter.HandleFunc("/tasks", ctrl.GetTasks).Methods("GET")
	authRouter.HandleFunc("/tasks/{id}", ctrl.GetTask).Methods("GET")
	authRouter.HandleFunc("/tasks", ctrl.CreateTask).Methods("POST")
	authRouter.HandleFunc("/tasks/{id}", ctrl.CompleteTask).Methods("PUT")
	authRouter.HandleFunc("/tasks/{id}", ctrl.DeleteTask).Methods("DELETE")
	authRouter.HandleFunc("/tasks/history/{id}", ctrl.DeleteTaskHistory).Methods("DELETE")

	handler := cors.Default().Handler(defaultRouter)
	log.Fatal(http.ListenAndServe(":8090", handler))
}
