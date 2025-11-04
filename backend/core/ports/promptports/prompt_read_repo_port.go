package promptports

import (
	"context"

	"ai_hub.com/app/core/domain/promptdomain"
	"ai_hub.com/app/core/domain/taskdomain"
)

type PromptReadRepository interface {
	FindByID(ctx context.Context, id promptdomain.PromptID) (*promptdomain.Prompt, error)

	FindAllByTask(ctx context.Context, taskID taskdomain.TaskID) ([]*promptdomain.Prompt, error)

	FindByIDs(ctx context.Context, ids []promptdomain.PromptID) ([]*promptdomain.Prompt, error)
}
