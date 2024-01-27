package router

import (
	"github.com/gorilla/mux"
	"github.com/jesper-nord/recurringly-backend/auth"
	"github.com/jesper-nord/recurringly-backend/task"
	"github.com/rs/cors"
	"log"
	"net/http"
	"os"
)

func New(taskService task.Service, authService auth.Service) http.Handler {
	authCtrl := auth.Controller{Service: authService}
	router := mux.NewRouter()
	router.HandleFunc("/api/login", authCtrl.Login).Methods("POST")
	router.HandleFunc("/api/register", authCtrl.Register).Methods("POST")
	router.HandleFunc("/api/refresh", authCtrl.RefreshToken).Methods("POST")

	authRouter := router.PathPrefix("/api/auth").Subrouter()
	authRouter.Use(JwtMiddleware)

	ctrl := task.Controller{Service: taskService}
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

	return c.Handler(router)
}
