package apikeyports

import "context"

type APIKeyBalancePort interface {
	SetBalance(ctx context.Context, apiKeyID string, newBalance float64, updatedByAdminID string) error
}
