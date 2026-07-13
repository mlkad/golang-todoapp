package tasks_postgres_repository

import (
	"time"

	"github.com/mlkad/golang-todoapp/internal/core/domain"
)

type TaskModel struct {
	ID          int
	Version     int
	Title       string
	Description *string
	Completed   bool
	CreatedAt   time.Time
	CompletedAt *time.Time
	AuthorID    int
}

func taskDomainsFromModels(taskModels []TaskModel) []domain.Task {
	domains := make([]domain.Task, len(taskModels))
	for i, model := range taskModels {
		domains[i] = taskDomainFromModel(model)
	}
	return domains
}

func taskDomainFromModel(taskModel TaskModel) domain.Task {
	return domain.NewTask(
		taskModel.ID,
		taskModel.Version,
		taskModel.Title,
		taskModel.Description,
		taskModel.Completed,
		taskModel.CreatedAt,
		taskModel.CompletedAt,
		taskModel.AuthorID,
	)
}
