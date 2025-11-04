package apikeyports

import "context"

type ModelInfo struct {
	CostPerToken float64

	MaxTokens int64
}

type APIKeyMetaPort interface {
	SetModelInfo(ctx context.Context, apiKeyID string, info ModelInfo) error
}
