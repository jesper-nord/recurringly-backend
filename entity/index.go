package entity

import "gorm.io/gorm"

type Task struct {
	gorm.Model
	Name     string
	Schedule int64
}
