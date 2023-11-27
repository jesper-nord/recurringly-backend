package dto

import "time"

type Task struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	History []TaskHistory
}

type TaskHistory struct {
	ID          string    `json:"id"`
	CompletedAt time.Time `json:"completed_at"`
}

type CreateTaskRequest struct {
	Name string `json:"name"`
}

type CompleteTaskRequest struct {
	ID string `json:"id"`
}
