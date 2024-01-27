package task

import (
	"errors"
	"github.com/jesper-nord/recurringly-backend/auth"
	"time"
)

type taskService struct {
	taskRepository Repository
}

func NewService(repository Repository) Service {
	return &taskService{
		taskRepository: repository,
	}
}

func (t *taskService) GetTasks(userId auth.UserId) ([]Task, error) {
	return t.taskRepository.FindAllTasks(userId)
}

func (t *taskService) GetTaskById(userId auth.UserId, taskId TaskId) (*Task, error) {
	return t.taskRepository.FindOneTask(userId, taskId)
}

func (t *taskService) CreateTask(userId auth.UserId, name string) (*Task, error) {
	if len(name) == 0 {
		return nil, errors.New("invalid task name")
	}

	task := Task{
		Name:   name,
		UserID: userId,
	}
	return t.taskRepository.SaveTask(&task)
}

func (t *taskService) EditTaskName(userId auth.UserId, taskId TaskId, name string) (*Task, error) {
	if len(name) == 0 {
		return nil, errors.New("invalid task name")
	}

	task, err := t.taskRepository.FindOneTask(userId, taskId)
	if err != nil {
		return &Task{}, err
	}
	task.Name = name
	return t.taskRepository.SaveTask(task)
}

func (t *taskService) CompleteTask(userId auth.UserId, taskId TaskId) (*Task, error) {
	return t.taskRepository.CompleteTask(userId, taskId)
}

func (t *taskService) DeleteTask(userId auth.UserId, taskId TaskId) error {
	task, err := t.taskRepository.FindOneTask(userId, taskId)
	if err != nil {
		return err
	}
	return t.taskRepository.DeleteTask(task)
}

func (t *taskService) EditTaskHistory(userId auth.UserId, taskId TaskId, taskHistoryId TaskHistoryId, completedAt time.Time) (*TaskHistory, error) {
	_, err := t.taskRepository.FindOneTask(userId, taskId)
	if err != nil {
		return nil, err
	}
	history, err := t.taskRepository.FindTaskHistory(taskId, taskHistoryId)
	if err != nil {
		return nil, err
	}

	history.CompletedAt = completedAt
	return t.taskRepository.SaveTaskHistory(history)
}

func (t *taskService) DeleteTaskHistory(userId auth.UserId, taskId TaskId, taskHistoryId TaskHistoryId) error {
	_, err := t.taskRepository.FindOneTask(userId, taskId)
	if err != nil {
		return err
	}
	history, err := t.taskRepository.FindTaskHistory(taskId, taskHistoryId)
	if err != nil {
		return err
	}

	return t.taskRepository.DeleteTaskHistory(history)
}
