package dto

import "time"

type Task struct {
	ID      string        `json:"id"`
	Name    string        `json:"name"`
	History []TaskHistory `json:"history"`
}

type TaskHistory struct {
	ID          string    `json:"id"`
	CompletedAt time.Time `json:"completed_at"`
}

type User struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

type CreateTaskRequest struct {
	Name string `json:"name"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	User   User      `json:"user"`
	Tokens TokenPair `json:"tokens"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenResponse struct {
	Tokens TokenPair `json:"tokens"`
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
