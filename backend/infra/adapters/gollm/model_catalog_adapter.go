package gollm

import (
	"context"

	"ai_hub.com/app/core/ports/apikeyports"
	"github.com/grokify/gollm"
)

type LLMGoModelCatalog struct{}

func NewLLMGoModelCatalog() *LLMGoModelCatalog { return &LLMGoModelCatalog{} }

func (a *LLMGoModelCatalog) GetModelInfoOrThrow(
	ctx context.Context,
	provider string,
	modelName string,
	keyValue string,
) (*apikeyports.ModelInfo, error) {
	providerName := gollm.ProviderName(provider)

	client, err := gollm.NewClient(gollm.ClientConfig{
		Provider: providerName,
		APIKey:   keyValue,
	})
	if err != nil {
		return nil, err
	}
	defer client.Close()

	mt := 1
	temp := 0.0

	_, err = client.CreateChatCompletion(ctx, &gollm.ChatCompletionRequest{
		Model:       modelName,
		Messages:    []gollm.Message{{Role: gollm.RoleUser, Content: "ping"}},
		MaxTokens:   &mt,
		Temperature: &temp,
	})
	if err != nil {
		return nil, err
	}

	return &apikeyports.ModelInfo{}, nil
}
