package adminports

import (
	"context"

	"ai_hub.com/app/core/domain/admindomain"
)

type AdminWriteRepository interface {
	Create(ctx context.Context, admin *admindomain.Admin) (*admindomain.Admin, error)

	Update(ctx context.Context, admin *admindomain.Admin) (*admindomain.Admin, error)

	Delete(ctx context.Context, id string) error
}
