package controller

import (
	"encoding/json"
	"gorm.io/gorm"
	"io"
	"net/http"
	"recurringly-backend/dto"
	"recurringly-backend/entity"
)

type TasksController struct {
	Database *gorm.DB
}

func (c TasksController) GetTasks(w http.ResponseWriter, r *http.Request) {
	var tasks []entity.Task
	c.Database.Find(&tasks)
	json.NewEncoder(w).Encode(&tasks)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (c TasksController) CreateTask(w http.ResponseWriter, r *http.Request) {
	var request dto.CreateTaskRequest
	body, _ := io.ReadAll(r.Body)
	err := json.Unmarshal(body, &request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	task := entity.Task{
		Name:     request.Name,
		Schedule: request.Schedule,
	}
	result := c.Database.Create(&task)
	if result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(&task)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}
