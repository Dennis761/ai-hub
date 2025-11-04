package taskports

import (
	"context"

	"ai_hub.com/app/core/domain/taskdomain"
)

type TaskWriteRepository interface {
	Create(ctx context.Context, task taskdomain.Task) (*taskdomain.Task, error)

	Update(ctx context.Context, task taskdomain.Task) (*taskdomain.Task, error)

	Delete(ctx context.Context, _id taskdomain.TaskID) error
}
