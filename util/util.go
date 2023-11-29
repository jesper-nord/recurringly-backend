package util

import (
	"recurringly-backend/dto"
	"recurringly-backend/entity"
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
