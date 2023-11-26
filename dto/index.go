package dto

type CreateTaskRequest struct {
	Name     string `json:"name"`
	Schedule int64  `json:"schedule"`
}

type CompleteTaskRequest struct {
	ID uint `json:"id"`
}
