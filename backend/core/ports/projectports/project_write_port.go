package projectports

import (
	"context"

	"ai_hub.com/app/core/domain/projectdomain"
)

type ProjectWriteRepository interface {
	Create(ctx context.Context, project *projectdomain.Project) (*projectdomain.Project, error)

	Update(ctx context.Context, project *projectdomain.Project) (*projectdomain.Project, error)

	Delete(ctx context.Context, id projectdomain.ProjectID) error
}
