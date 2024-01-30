package task

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/jesper-nord/recurringly-backend/auth"
	"io"
	"log"
	"net/http"
	"sort"
	"strconv"
)

type Controller struct {
	Service Service
}

func (c Controller) GetTasks(w http.ResponseWriter, r *http.Request) {
	userId := getUserIdFromRequest(r)
	tasks, err := c.Service.GetTasks(userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tasksToApiModel(tasks))
}

func (c Controller) GetTask(w http.ResponseWriter, r *http.Request) {
	userId := getUserIdFromRequest(r)
	taskId := toTaskId(mux.Vars(r)["id"])

	task, err := c.Service.GetTaskById(userId, taskId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(taskToApiModel(task))
}

func (c Controller) CreateTask(w http.ResponseWriter, r *http.Request) {
	userId := getUserIdFromRequest(r)
	var request CreateTaskRequest
	body, _ := io.ReadAll(r.Body)
	err := json.Unmarshal(body, &request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	createdTask, err := c.Service.CreateTask(userId, request.Name)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Printf("created task: %d for user: %d", createdTask.ID, userId)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(taskToApiModel(createdTask))
}

func (c Controller) EditTask(w http.ResponseWriter, r *http.Request) {
	userId := getUserIdFromRequest(r)
	taskId := toTaskId(mux.Vars(r)["id"])

	var request CreateTaskRequest
	body, _ := io.ReadAll(r.Body)
	err := json.Unmarshal(body, &request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	task, err := c.Service.EditTaskName(userId, taskId, request.Name)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Printf("updated task: %d for user: %d", task.ID, userId)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(taskToApiModel(task))
}

func (c Controller) CompleteTask(w http.ResponseWriter, r *http.Request) {
	userId := getUserIdFromRequest(r)
	taskId := toTaskId(mux.Vars(r)["id"])

	task, err := c.Service.CompleteTask(userId, taskId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	log.Printf("completed task: %d for user: %d", task.ID, userId)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(taskToApiModel(task))
}

func (c Controller) DeleteTask(w http.ResponseWriter, r *http.Request) {
	userId := getUserIdFromRequest(r)
	taskId := toTaskId(mux.Vars(r)["id"])

	err := c.Service.DeleteTask(userId, taskId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	log.Printf("deleted task: %d for user: %d", taskId, userId)
	w.WriteHeader(http.StatusOK)
}

func (c Controller) EditTaskHistory(w http.ResponseWriter, r *http.Request) {
	userId := getUserIdFromRequest(r)
	taskId := toTaskId(mux.Vars(r)["id"])
	taskHistoryId := toTaskHistoryId(mux.Vars(r)["historyId"])

	var request EditTaskHistoryRequest
	body, _ := io.ReadAll(r.Body)
	err := json.Unmarshal(body, &request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = c.Service.EditTaskHistory(userId, taskId, taskHistoryId, request.CompletedAt)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	log.Printf("updated task history: %d in task: %d for user: %d", taskHistoryId, taskId, userId)
	w.WriteHeader(http.StatusOK)
}

func (c Controller) DeleteTaskHistory(w http.ResponseWriter, r *http.Request) {
	userId := getUserIdFromRequest(r)
	taskId := toTaskId(mux.Vars(r)["id"])
	taskHistoryId := toTaskHistoryId(mux.Vars(r)["historyId"])

	err := c.Service.DeleteTaskHistory(userId, taskId, taskHistoryId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	log.Printf("deleted task history: %d in task: %d for user: %d", taskHistoryId, taskId, userId)
	w.WriteHeader(http.StatusOK)
}

func getUserIdFromRequest(r *http.Request) auth.UserId {
	id, _ := strconv.Atoi(r.Context().Value("user").(string))
	return auth.UserId(id)
}

func toTaskId(s string) TaskId {
	return TaskId(parseUint(s))
}

func toTaskHistoryId(s string) TaskHistoryId {
	return TaskHistoryId(parseUint(s))
}

func parseUint(s string) uint64 {
	parsed, _ := strconv.ParseUint(s, 10, 32)
	return parsed
}

func tasksToApiModel(tasks []Task) []ApiTask {
	result := make([]ApiTask, 0)
	for _, task := range tasks {
		result = append(result, taskToApiModel(&task))
	}
	return result
}

func taskToApiModel(task *Task) ApiTask {
	return ApiTask{
		ID:      task.ID,
		Name:    task.Name,
		History: taskHistoriesToApiModel(task.History),
	}
}

func taskHistoriesToApiModel(history []TaskHistory) []ApiTaskHistory {
	result := make([]ApiTaskHistory, 0)
	for _, th := range history {
		result = append(result, taskHistoryToApiModel(th))
	}
	// sort by completion date, descending
	sort.Slice(result, func(i, j int) bool {
		return result[i].CompletedAt.After(result[j].CompletedAt)
	})
	return result
}

func taskHistoryToApiModel(history TaskHistory) ApiTaskHistory {
	return ApiTaskHistory{
		ID:          history.ID,
		CompletedAt: history.CompletedAt,
	}
}
