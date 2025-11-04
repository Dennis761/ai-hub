package apikeyports

import (
	"context"

	"ai_hub.com/app/core/domain/apikeydomain"
)

type APIKeyWriteRepository interface {
	Create(ctx context.Context, entity *apikeydomain.APIKey) (*apikeydomain.APIKey, error)

	Update(ctx context.Context, entity *apikeydomain.APIKey) (*apikeydomain.APIKey, error)

	Delete(ctx context.Context, id apikeydomain.APIKeyID) error
}
