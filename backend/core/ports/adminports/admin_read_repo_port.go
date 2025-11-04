package adminports

import (
	"context"

	"ai_hub.com/app/core/domain/admindomain"
)

type AdminReadRepository interface {
	FindByID(ctx context.Context, id string) (*admindomain.Admin, error)

	FindByEmail(ctx context.Context, email string) (*admindomain.Admin, error)
}
