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
	taskId, err := uuid.Parse(params["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	task, err := c.getTask(taskId)
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
	params := mux.Vars(r)
	taskId, err := uuid.Parse(params["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	task, err := c.getTask(taskId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	err = c.Database.Model(&task).Association("History").Append(&entity.TaskHistory{
		CompletedAt: time.Now(),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(util.TaskToApiModel(task))
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

func (c TasksController) getTask(taskId uuid.UUID) (entity.Task, error) {
	var task entity.Task
	err := c.Database.Model(&entity.Task{}).Preload("History").Find(&task, taskId).Error
	return task, err
}
