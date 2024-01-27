package task

import (
	"time"
)

type ApiTask struct {
	ID      string           `json:"id"`
	Name    string           `json:"name"`
	History []ApiTaskHistory `json:"history"`
}

type ApiTaskHistory struct {
	ID          string    `json:"id"`
	CompletedAt time.Time `json:"completed_at"`
}

type CreateTaskRequest struct {
	Name string `json:"name"`
}

type EditTaskHistoryRequest struct {
	CompletedAt time.Time `json:"completed_at"`
}
