package promptports

import "context"

type BillingRunResult struct {
	ResponseText     string
	Cost             float64
	BalanceRemaining float64
	ModelName        string
	Provider         string
}

type RunWithBillingParams struct {
	ModelID         string
	FormattedPrompt string
	RunByAdminID    string
}

type LLMBillingRunnerPort interface {
	RunWithBilling(ctx context.Context, params RunWithBillingParams) (*BillingRunResult, error)
}
