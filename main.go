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
	db.AutoMigrate(&entity.Task{})

	ctrl := controller.TasksController{Database: db}
	router.HandleFunc("/tasks", ctrl.GetTasks).Methods("GET")
	router.HandleFunc("/tasks", ctrl.CreateTask).Methods("POST")

	handler := cors.Default().Handler(router)
	log.Fatal(http.ListenAndServe(":8090", handler))
}
