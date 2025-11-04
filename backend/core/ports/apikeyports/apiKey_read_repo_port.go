package apikeyports

import (
	"context"

	"ai_hub.com/app/core/domain/apikeydomain"
)

type APIKeyReadRepository interface {
	FindAllByOwner(ctx context.Context, ownerID apikeydomain.OwnerID) ([]*apikeydomain.APIKey, error)

	FindByID(ctx context.Context, id apikeydomain.APIKeyID) (*apikeydomain.APIKey, error)

	FindAllByKeyName(ctx context.Context, keyName apikeydomain.APIKeyName) ([]*apikeydomain.APIKey, error)

	GetModelByID(ctx context.Context, modelID string) (*apikeydomain.APIKey, error)
}
