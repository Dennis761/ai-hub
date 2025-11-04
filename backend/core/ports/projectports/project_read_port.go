package projectports

import (
	"context"

	"ai_hub.com/app/core/domain/projectdomain"
)

type ProjectReadRepository interface {
	FindByID(ctx context.Context, id projectdomain.ProjectID) (*projectdomain.Project, error)

	FindByName(ctx context.Context, name projectdomain.ProjectName) (*projectdomain.Project, error)

	FindAllAccessibleByAdmin(ctx context.Context, adminID string) ([]*projectdomain.Project, error)
}
