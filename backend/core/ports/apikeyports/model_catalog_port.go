package apikeyports

import "context"

type ModelCatalogPort interface {
	GetModelInfoOrThrow(ctx context.Context, provider string, modelName string, keyValue string) (*ModelInfo, error)
}
