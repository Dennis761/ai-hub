package taskports

import (
	"context"

	"ai_hub.com/app/core/domain/taskdomain"
)

type TaskReadRepository interface {
	FindByID(ctx context.Context, _id taskdomain.TaskID) (*taskdomain.Task, error)

	FindAllByProject(ctx context.Context, projectID taskdomain.TaskProjectID) ([]taskdomain.Task, error)
}
