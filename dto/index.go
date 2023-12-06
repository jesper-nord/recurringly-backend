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

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type AuthResponse struct {
	Tokens TokenPair `json:"tokens"`
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
