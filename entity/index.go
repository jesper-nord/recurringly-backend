package entity

import (
	"gorm.io/gorm"
	"time"
)

type Task struct {
	gorm.Model
	Name    string
	History []TaskHistory
}

type TaskHistory struct {
	gorm.Model
	CompletedAt time.Time
	TaskID      uint
}
