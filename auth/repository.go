package auth

import (
	"errors"
	"github.com/jesper-nord/recurringly-backend/domain"
	"gorm.io/gorm"
)

type User struct {
	domain.Model
	Username string
	Password string
}

type authRepository struct {
	Database *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &authRepository{
		Database: db,
	}
}

func (a *authRepository) GetUser(userId UserId) (*User, error) {
	var user User
	err := a.Database.Take(&user, userId).Error
	return &user, err
}

func (a *authRepository) FindUser(username string) (*User, error) {
	var user User
	err := a.Database.Where("username = ?", username).Take(&user).Error
	return &user, err
}

func (a *authRepository) SaveUser(user *User) (*User, error) {
	return user, a.Database.Save(user).Error
}

func (a *authRepository) Migrate() error {
	return errors.Join(a.Database.AutoMigrate(&User{}))
}
