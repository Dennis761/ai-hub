package gollm

import (
	"context"
	"fmt"

	apikeydomain "ai_hub.com/app/core/domain/apikeydomain"
	apikeyports "ai_hub.com/app/core/ports/apikeyports"
	promptports "ai_hub.com/app/core/ports/promptports"
)

type LLMBillingRunner struct {
	apiKeys  apikeyports.APIKeyReadRepository
	crypto   apikeyports.CryptoPort
	llm      LLMClient
	balances apikeyports.APIKeyBalancePort
}

func NewLLMBillingRunner(
	apiKeys apikeyports.APIKeyReadRepository,
	crypto apikeyports.CryptoPort,
	llm LLMClient,
	balances apikeyports.APIKeyBalancePort,
) *LLMBillingRunner {
	return &LLMBillingRunner{
		apiKeys:  apiKeys,
		crypto:   crypto,
		llm:      llm,
		balances: balances,
	}
}

// RunWithBilling:
// -receives a prompt,
// -calls the model via gollm,
// -calculates the cost,
// -updates the balance,
// -returns the result to the domain port.
func (r *LLMBillingRunner) RunWithBilling(
	ctx context.Context,
	params promptports.RunWithBillingParams,
) (*promptports.BillingRunResult, error) {

	keyID, err := apikeydomain.NewAPIKeyID(params.ModelID)
	if err != nil {
		return nil, err
	}

	entity, err := r.apiKeys.FindByID(ctx, keyID)
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return nil, fmt.Errorf("api key not found: %s", params.ModelID)
	}

	p := entity.ToPrimitives()

	apiKey, err := r.crypto.Decrypt(p.KeyValue)
	if err != nil {
		return nil, fmt.Errorf("decrypt api key: %w", err)
	}

	responseText, cost, err := r.llm.Run(
		ctx,
		p.Provider,
		p.ModelName,
		apiKey,
		params.FormattedPrompt,
	)
	if err != nil {
		return nil, err
	}

	newBalance := p.Balance - cost

	if err := r.balances.SetBalance(ctx, p.ID, newBalance, params.RunByAdminID); err != nil {
		return nil, err
	}

	return &promptports.BillingRunResult{
		ResponseText:     responseText,
		Cost:             cost,
		BalanceRemaining: newBalance,
		ModelName:        p.ModelName,
		Provider:         p.Provider,
	}, nil
}
