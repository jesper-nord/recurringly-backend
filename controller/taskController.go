package controller

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/jesper-nord/recurringly-backend/dto"
	"github.com/jesper-nord/recurringly-backend/entity"
	"github.com/jesper-nord/recurringly-backend/util"
	"gorm.io/gorm"
	"io"
	"net/http"
	"time"
)

type TaskController struct {
	Database *gorm.DB
}

func (c TaskController) GetTasks(w http.ResponseWriter, r *http.Request) {
	userId := util.GetUserIdFromRequest(r)
	var tasks []entity.Task
	err := c.Database.Where("user_id = ?", userId).Preload("History").Find(&tasks).Error
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(util.TasksToApiModel(tasks))
}

func (c TaskController) GetTask(w http.ResponseWriter, r *http.Request) {
	userId := util.GetUserIdFromRequest(r)
	params := mux.Vars(r)
	taskId := params["id"]

	task, err := c.getTask(taskId, userId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(util.TaskToApiModel(task))
}

func (c TaskController) CreateTask(w http.ResponseWriter, r *http.Request) {
	userId, err := uuid.Parse(util.GetUserIdFromRequest(r))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var request dto.CreateTaskRequest
	body, _ := io.ReadAll(r.Body)
	err = json.Unmarshal(body, &request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(request.Name) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	task := entity.Task{
		Name:   request.Name,
		UserID: userId,
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

func (c TaskController) EditTask(w http.ResponseWriter, r *http.Request) {
	userId := util.GetUserIdFromRequest(r)
	params := mux.Vars(r)
	taskId := params["id"]

	var request dto.CreateTaskRequest
	body, _ := io.ReadAll(r.Body)
	err := json.Unmarshal(body, &request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(request.Name) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	task, err := c.getTask(taskId, userId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	task.Name = request.Name
	err = c.Database.Save(task).Error
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(util.TaskToApiModel(task))
}

func (c TaskController) CompleteTask(w http.ResponseWriter, r *http.Request) {
	userId := util.GetUserIdFromRequest(r)
	params := mux.Vars(r)
	taskId := params["id"]

	task, err := c.getTask(taskId, userId)
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

func (c TaskController) DeleteTask(w http.ResponseWriter, r *http.Request) {
	userId := util.GetUserIdFromRequest(r)
	params := mux.Vars(r)
	c.Database.Where("id = ? AND user_id = ?", params["id"], userId).Delete(&entity.Task{})
	w.WriteHeader(http.StatusOK)
}

func (c TaskController) EditTaskHistory(w http.ResponseWriter, r *http.Request) {
	userId := util.GetUserIdFromRequest(r)
	params := mux.Vars(r)
	taskId := params["id"]
	taskHistoryId := params["historyId"]

	var request dto.EditTaskHistoryRequest
	body, _ := io.ReadAll(r.Body)
	err := json.Unmarshal(body, &request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = c.getTask(taskId, userId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err = c.Database.Model(&entity.TaskHistory{}).Where("id = ? AND task_id = ?", taskHistoryId, taskId).Update("completed_at", request.CompletedAt).Error
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (c TaskController) DeleteTaskHistory(w http.ResponseWriter, r *http.Request) {
	userId := util.GetUserIdFromRequest(r)
	params := mux.Vars(r)
	taskId := params["id"]
	taskHistoryId := params["historyId"]

	var task entity.Task
	err := c.Database.Where("id = ? AND user_id = ?", taskId, userId).Take(&task).Error
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err = c.Database.Where("id = ? AND task_id = ?", taskHistoryId, taskId).Delete(&entity.TaskHistory{}).Error
	w.WriteHeader(http.StatusOK)
}

func (c TaskController) getTask(taskId string, userId string) (entity.Task, error) {
	var task entity.Task
	err := c.Database.Where("id = ? AND user_id = ?", taskId, userId).Preload("History").Take(&task).Error
	return task, err
}
