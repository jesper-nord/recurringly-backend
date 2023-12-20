package util

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jesper-nord/recurringly-backend/dto"
	"github.com/jesper-nord/recurringly-backend/entity"
	"net/http"
	"os"
	"sort"
)

func TasksToApiModel(tasks []entity.Task) []dto.Task {
	result := make([]dto.Task, 0)
	for _, task := range tasks {
		result = append(result, TaskToApiModel(task))
	}
	return result
}

func TaskToApiModel(task entity.Task) dto.Task {
	return dto.Task{
		ID:      task.ID.String(),
		Name:    task.Name,
		History: taskHistoriesToApiModel(task.History),
	}
}

func taskHistoriesToApiModel(history []entity.TaskHistory) []dto.TaskHistory {
	result := make([]dto.TaskHistory, 0)
	for _, th := range history {
		result = append(result, TaskHistoryToApiModel(th))
	}
	// sort by completion date, descending
	sort.Slice(result, func(i, j int) bool {
		return result[i].CompletedAt.After(result[j].CompletedAt)
	})
	return result
}

func TaskHistoryToApiModel(history entity.TaskHistory) dto.TaskHistory {
	return dto.TaskHistory{
		ID:          history.ID.String(),
		CompletedAt: history.CompletedAt,
	}
}

func ParseJwt(token string, claims jwt.Claims) (*jwt.Token, error) {
	return jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SIGNING_SECRET")), nil
	})
}

func GetUserIdFromRequest(r *http.Request) string {
	return r.Context().Value("user").(string)
}
