package gollm

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	gollmsdk "github.com/grokify/gollm"
)

type LLMClient interface {
	Run(ctx context.Context, provider, model, apiKey, prompt string) (response string, cost float64, err error)
}

type GollmClientAdapter struct{}

func NewGollmClientAdapter() *GollmClientAdapter { return &GollmClientAdapter{} }

//
// PRICING CACHE
//

// priceItem describes a single model entry from https://www.llm-prices.com/current-v1.json
// Example entry (simplified):
//
//	{
//	  "id": "gemini-2.0-flash",
//	  "vendor": "google",
//	  "name": "Gemini 2.0 Flash",
//	  "input": 0.035,
//	  "output": 0.14
//	}
//
// Notes:
//   - "id" is the machine name for the model (sometimes matches exactly what we call the model).
//   - "name" is a human-friendly label ("Gemini 2.0 Flash"), which may differ in formatting
//     but represents the same thing.
//   - "input"/"output" are USD per 1M tokens, not per 1 token.
type priceItem struct {
	ID     string   `json:"id"`
	Vendor string   `json:"vendor"`
	Name   string   `json:"name"`
	Input  *float64 `json:"input"`  // USD per 1M input tokens
	Output *float64 `json:"output"` // USD per 1M output tokens
}

// pricesPayload is the top-level structure returned by llm-prices.com
//
//	{
//	  "updated_at": "...",
//	  "prices": [ { ...priceItem... }, ... ]
//	}
type pricesPayload struct {
	UpdatedAt string      `json:"updated_at"`
	Prices    []priceItem `json:"prices"`
}

// We only want to fetch pricing data once per process lifetime.
// sync.Once guarantees that, and we keep the result in-memory.
var (
	pricesOnce  sync.Once
	pricesCache pricesPayload
	pricesErr   error
)

// loadPricesOnce downloads and caches the pricing data from llm-prices.com.
// After the first call, subsequent calls just return the cached data.

func loadPricesOnce() (pricesPayload, error) {
	pricesOnce.Do(func() {
		resp, err := http.Get("https://www.llm-prices.com/current-v1.json")
		if err != nil {
			pricesErr = fmt.Errorf("fetch llm prices failed: %w", err)
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			pricesErr = fmt.Errorf("read llm prices failed: %w", err)
			return
		}

		if err := json.Unmarshal(body, &pricesCache); err != nil {
			pricesErr = fmt.Errorf("unmarshal llm prices failed: %w", err)
			return
		}
	})
	return pricesCache, pricesErr
}

//
// PRICE LOOKUP LOGIC
//

// normalizePriceName converts the human-readable model name from the pricing JSON
// into the normalized kebab-case style we use internally.
//
// Examples:
//
//	"Gemini 2.0 Flash"   -> "gemini-2.0-flash"
//	"GPT-4o Mini"        -> "gpt-4o-mini"
//	"Claude 3.5 Sonnet"  -> "claude-3.5-sonnet"
//
// We ONLY normalize the provider's "Name", not our input model string.
// Why?
//   - Our own model string is already normalized before it reaches this code
//     (e.g. "gemini-2.0-flash").
//   - But the pricing file might have "Gemini 2.0 Flash" with spaces and caps.
//     Converting the pricing name to kebab-case lets us match them.
func normalizePriceName(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	s = strings.ReplaceAll(s, " ", "-")
	return s
}

// getPerTokenUSD returns the price per single token (input and output) for the given model.
// We intentionally do NOT require provider/vendor to match.
//
// Why we don't require vendor:
//   - Real-world quirk: you might call provider="gemini",
//     but the pricing file lists vendor="google" for those Gemini models.
//     If we forced vendor to match exactly, we'd fail to find a valid price
//     even though it's obviously the same model.
//
// How we match the model name:
//
//  1. First, try to match request `model` against `p.ID` exactly (case-insensitive).
//     Example: you pass "gpt-4o-mini", pricing has ID "gpt-4o-mini" → perfect match.
//
//  2. If that fails, try to match request `model` against normalized `p.Name`.
//     Example: you pass "gemini-2.0-flash",
//     pricing has Name "Gemini 2.0 Flash".
//     We normalize "Gemini 2.0 Flash" -> "gemini-2.0-flash",
//     and now we can match it.
//
// We only need the first match because we assume you're asking about one specific model.
func getPerTokenUSD(_provider, model string) (float64, float64, error) {
	cache, err := loadPricesOnce()
	if err != nil {
		return 0, 0, err
	}

	modelLower := strings.ToLower(model)

	var chosen *priceItem

	for _, p := range cache.Prices {
		idLower := strings.ToLower(p.ID)
		nameNorm := normalizePriceName(p.Name)

		// 1. strict ID match ("gpt-4o-mini" == "gpt-4o-mini")
		if idLower == modelLower {
			chosen = &p
			break
		}

		// 2. fallback match via normalized Name
		//    e.g. model "gemini-2.0-flash" vs Name "Gemini 2.0 Flash"
		if nameNorm == modelLower && chosen == nil {
			tmp := p
			chosen = &tmp
		}
	}

	if chosen == nil {
		// We didn't find any pricing that looks like this model.
		// This often means the pricing JSON uses a slightly different ID,
		// like "google-gemini-2.0-flash" instead of just "gemini-2.0-flash".
		return 0, 0, fmt.Errorf("model not found in pricing: %s", model)
	}

	if chosen.Input == nil || chosen.Output == nil {
		// Some entries might be missing pricing fields for edge cases.
		return 0, 0, fmt.Errorf("no price data for %s", model)
	}

	// The pricing file gives $ per 1,000,000 tokens.
	// We want $ per 1 token.
	inPerToken := *chosen.Input / 1_000_000.0
	outPerToken := *chosen.Output / 1_000_000.0

	return inPerToken, outPerToken, nil
}

func (a *GollmClientAdapter) Run(
	ctx context.Context,
	provider, model, apiKey, prompt string,
) (string, float64, error) {

	client, err := gollmsdk.NewClient(gollmsdk.ClientConfig{
		Provider: gollmsdk.ProviderName(provider),
		APIKey:   apiKey,
	})
	if err != nil {
		return "", 0, fmt.Errorf("provider init failed (%s): %w", provider, err)
	}
	defer client.Close()

	maxTokens := 256

	req := &gollmsdk.ChatCompletionRequest{
		Model: model,
		Messages: []gollmsdk.Message{
			{
				Role:    gollmsdk.RoleUser,
				Content: prompt,
			},
		},
		MaxTokens: &maxTokens,
	}

	start := time.Now()
	resp, err := client.CreateChatCompletion(ctx, req)
	_ = start
	if err != nil {
		return "", 0, err
	}

	if len(resp.Choices) == 0 {
		return "", 0, fmt.Errorf("empty response from provider=%s model=%s", provider, model)
	}

	answer := resp.Choices[0].Message.Content

	inPerTokUSD, outPerTokUSD, err := getPerTokenUSD(provider, model)
	if err != nil {
		return answer, 0, fmt.Errorf("price lookup failed: %w", err)
	}

	inputTokens := float64(resp.Usage.PromptTokens)
	outputTokens := float64(resp.Usage.CompletionTokens)

	costUSD := inputTokens*inPerTokUSD + outputTokens*outPerTokUSD

	return answer, costUSD, nil
}
