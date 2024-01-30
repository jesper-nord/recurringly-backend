package task

import (
	"github.com/jesper-nord/recurringly-backend/auth"
	"gorm.io/gorm"
	"time"
)

type Task struct {
	gorm.Model
	Name    string
	UserID  uint
	History []TaskHistory
}

type TaskHistory struct {
	gorm.Model
	CompletedAt time.Time
	TaskID      uint
}

type taskRepository struct {
	Database *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &taskRepository{
		Database: db,
	}
}

func (t *taskRepository) FindOneTask(userId auth.UserId, taskId TaskId) (*Task, error) {
	var task Task
	err := t.Database.Where("id = ? AND user_id = ?", taskId, userId).Preload("History").Take(&task).Error
	return &task, err
}

func (t *taskRepository) FindAllTasks(user auth.UserId) ([]Task, error) {
	var tasks []Task
	err := t.Database.Where("user_id = ?", user).Preload("History").Find(&tasks).Error
	return tasks, err
}

func (t *taskRepository) SaveTask(task *Task) (*Task, error) {
	return task, t.Database.Save(task).Error
}

func (t *taskRepository) DeleteTask(task *Task) error {
	return t.Database.Delete(&task).Error
}

func (t *taskRepository) CompleteTask(userId auth.UserId, taskId TaskId) (*Task, error) {
	task, err := t.FindOneTask(userId, taskId)
	if err != nil {
		return nil, err
	}
	return task, t.Database.Model(task).Association("History").Append(&TaskHistory{
		CompletedAt: time.Now(),
	})
}

func (t *taskRepository) FindTaskHistory(taskId TaskId, taskHistoryId TaskHistoryId) (*TaskHistory, error) {
	var taskHistory TaskHistory
	err := t.Database.Where("id = ? AND task_id = ?", taskHistoryId, taskId).Take(&taskHistory).Error
	return &taskHistory, err
}

func (t *taskRepository) SaveTaskHistory(taskHistory *TaskHistory) (*TaskHistory, error) {
	return taskHistory, t.Database.Save(taskHistory).Error
}

func (t *taskRepository) DeleteTaskHistory(taskHistory *TaskHistory) error {
	return t.Database.Delete(&taskHistory).Error
}

func (t *taskRepository) Migrate() error {
	return t.Database.AutoMigrate(&Task{}, &TaskHistory{})
}
