package task

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"sort"
)

type Controller struct {
	Service Service
}

func (c Controller) GetTasks(w http.ResponseWriter, r *http.Request) {
	userId := getUserUuidFromRequest(r)
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
	userId := getUserUuidFromRequest(r)
	taskId := getId(mux.Vars(r))

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
	userId := getUserUuidFromRequest(r)
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

	log.Printf("created task: '%s' for user '%s'", createdTask.ID.String(), userId)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(taskToApiModel(createdTask))
}

func (c Controller) EditTask(w http.ResponseWriter, r *http.Request) {
	userId := getUserUuidFromRequest(r)
	taskId := getId(mux.Vars(r))

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

	log.Printf("updated task: '%s' for user '%s'", task.ID.String(), userId)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(taskToApiModel(task))
}

func (c Controller) CompleteTask(w http.ResponseWriter, r *http.Request) {
	userId := getUserUuidFromRequest(r)
	taskId := getId(mux.Vars(r))

	task, err := c.Service.CompleteTask(userId, taskId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	log.Printf("completed task: '%s' for user '%s'", task.ID.String(), userId)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(taskToApiModel(task))
}

func (c Controller) DeleteTask(w http.ResponseWriter, r *http.Request) {
	userId := getUserUuidFromRequest(r)
	taskId := getId(mux.Vars(r))

	err := c.Service.DeleteTask(userId, taskId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	log.Printf("deleted task: '%s' for user '%s'", taskId, userId)
	w.WriteHeader(http.StatusOK)
}

func (c Controller) EditTaskHistory(w http.ResponseWriter, r *http.Request) {
	userId := getUserUuidFromRequest(r)
	params := mux.Vars(r)
	taskId := getId(params)
	taskHistoryId := getUuid("historyId", params)

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

	log.Printf("updated task history: '%s' in task '%s' for user '%s'", taskHistoryId, taskId, userId)
	w.WriteHeader(http.StatusOK)
}

func (c Controller) DeleteTaskHistory(w http.ResponseWriter, r *http.Request) {
	userId := getUserUuidFromRequest(r)
	params := mux.Vars(r)
	taskId := getId(params)
	taskHistoryId := getUuid("historyId", params)

	err := c.Service.DeleteTaskHistory(userId, taskId, taskHistoryId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	log.Printf("deleted task history: '%s' in task '%s' for user '%s'", taskHistoryId, taskId, userId)
	w.WriteHeader(http.StatusOK)
}

func getUserUuidFromRequest(r *http.Request) uuid.UUID {
	parsed, _ := uuid.Parse(r.Context().Value("user").(string))
	return parsed
}

func getId(params map[string]string) uuid.UUID {
	return getUuid("id", params)
}

func getUuid(name string, params map[string]string) uuid.UUID {
	id, _ := uuid.Parse(params[name])
	return id
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
		ID:      task.ID.String(),
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
		ID:          history.ID.String(),
		CompletedAt: history.CompletedAt,
	}
}
