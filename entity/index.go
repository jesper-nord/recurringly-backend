package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Model struct {
	ID        uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Task struct {
	Model
	Name    string
	History []TaskHistory
}

type TaskHistory struct {
	Model
	CompletedAt time.Time
	TaskID      uuid.UUID
}

type User struct {
	Model
	Email    string
	Password string
}
