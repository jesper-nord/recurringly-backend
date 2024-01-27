package task

import (
	"github.com/google/uuid"
	"github.com/jesper-nord/recurringly-backend/auth"
	"time"
)

type TaskId = uuid.UUID
type TaskHistoryId = uuid.UUID

type Service interface {
	GetTasks(userId auth.UserId) ([]Task, error)
	GetTaskById(userId auth.UserId, taskId TaskId) (*Task, error)
	CreateTask(userId auth.UserId, name string) (*Task, error)
	EditTaskName(userId auth.UserId, taskId TaskId, name string) (*Task, error)
	CompleteTask(userId auth.UserId, taskId TaskId) (*Task, error)
	DeleteTask(userId auth.UserId, taskId TaskId) error
	EditTaskHistory(userId auth.UserId, taskId TaskId, taskHistoryId TaskHistoryId, completedAt time.Time) (*TaskHistory, error)
	DeleteTaskHistory(userId auth.UserId, taskId TaskId, taskHistoryId TaskHistoryId) error
}

type Repository interface {
	FindAllTasks(userId auth.UserId) ([]Task, error)
	FindOneTask(userId auth.UserId, taskId TaskId) (*Task, error)
	SaveTask(task *Task) (*Task, error)
	DeleteTask(task *Task) error
	CompleteTask(userId auth.UserId, taskId TaskId) (*Task, error)
	FindTaskHistory(taskId TaskId, taskHistoryId TaskHistoryId) (*TaskHistory, error)
	SaveTaskHistory(taskHistory *TaskHistory) (*TaskHistory, error)
	DeleteTaskHistory(taskHistory *TaskHistory) error
	Migrate() error
}
