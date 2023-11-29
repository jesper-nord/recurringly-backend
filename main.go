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
	router := mux.NewRouter()

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

	ctrl := controller.TaskController{Database: db}
	router.HandleFunc("/tasks", ctrl.GetTasks).Methods("GET")
	router.HandleFunc("/tasks/{id}", ctrl.GetTask).Methods("GET")
	router.HandleFunc("/tasks", ctrl.CreateTask).Methods("POST")
	router.HandleFunc("/tasks/{id}", ctrl.CompleteTask).Methods("PUT")
	router.HandleFunc("/tasks/{id}", ctrl.DeleteTask).Methods("DELETE")
	router.HandleFunc("/tasks/history/{id}", ctrl.DeleteTaskHistory).Methods("DELETE")
	authCtrl := controller.AuthController{Database: db}
	router.HandleFunc("/login", authCtrl.Login).Methods("POST")
	router.HandleFunc("/register", authCtrl.Register).Methods("POST")

	handler := cors.Default().Handler(router)
	log.Fatal(http.ListenAndServe(":8090", handler))
}
