package controller

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"io"
	"net/http"
	"recurringly-backend/dto"
	"recurringly-backend/entity"
	"recurringly-backend/util"
	"time"
)

type TasksController struct {
	Database *gorm.DB
}

func (c TasksController) GetTasks(w http.ResponseWriter, r *http.Request) {
	var tasks []entity.Task
	err := c.Database.Model(&entity.Task{}).Preload("History").Find(&tasks).Error
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(util.TasksToApiModel(tasks))
}

func (c TasksController) GetTask(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	var task entity.Task
	err := c.Database.Model(&entity.Task{}).Preload("History").Find(&task, id).Error
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(util.TaskToApiModel(task))
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
		Name: request.Name,
	}
	err = c.Database.Create(&task).Error
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(util.TaskToApiModel(task))
}

func (c TasksController) CompleteTask(w http.ResponseWriter, r *http.Request) {
	var request dto.CompleteTaskRequest
	body, _ := io.ReadAll(r.Body)
	err := json.Unmarshal(body, &request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	taskId, err := uuid.Parse(request.ID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	taskHistoryEntry := entity.TaskHistory{
		TaskID:      taskId,
		CompletedAt: time.Now(),
	}
	err = c.Database.Create(&taskHistoryEntry).Error
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(util.TaskHistoryToApiModel(taskHistoryEntry))
}

func (c TasksController) DeleteTask(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	taskId, err := uuid.Parse(params["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	c.Database.Delete(&entity.Task{}, taskId)
	w.WriteHeader(http.StatusOK)
}
