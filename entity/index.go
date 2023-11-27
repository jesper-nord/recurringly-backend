package entity

import (
	"database/sql"
	"github.com/google/uuid"
	"time"
)

type Model struct {
	ID        uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt sql.NullTime `gorm:"index"`
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
