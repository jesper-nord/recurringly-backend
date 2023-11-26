package entity

import (
	"gorm.io/gorm"
	"time"
)

type Task struct {
	gorm.Model
	Name     string
	Schedule int64
	History  []TaskHistory
}

type TaskHistory struct {
	gorm.Model
	DoneAt time.Time
	TaskID uint
}
