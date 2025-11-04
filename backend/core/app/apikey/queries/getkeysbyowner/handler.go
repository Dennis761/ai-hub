package apikeyapp

import (
	"context"

	"ai_hub.com/app/core/domain/apikeydomain"
	"ai_hub.com/app/core/ports/apikeyports"
)

type GetKeysByOwnerHandler struct {
	readRepo   apikeyports.APIKeyReadRepository
	cryptoPort apikeyports.CryptoPort
}

func NewGetKeysByOwnerHandler(
	readRepo apikeyports.APIKeyReadRepository,
	crypto apikeyports.CryptoPort,
) *GetKeysByOwnerHandler {
	return &GetKeysByOwnerHandler{
		readRepo:   readRepo,
		cryptoPort: crypto,
	}
}

func (h *GetKeysByOwnerHandler) GetKeysByOwner(
	ctx context.Context,
	q GetKeysByOwnerQuery,
) (GetKeysByOwnerResult, error) {
	// normalize owner id
	ownerID, err := apikeydomain.NewOwnerID(q.OwnerID)
	if err != nil {
		return nil, err
	}

	// load all keys for this owner
	apiKeys, err := h.readRepo.FindAllByOwner(ctx, ownerID)
	if err != nil {
		return nil, err
	}

	// map domain aggregates to DTO
	out := make(GetKeysByOwnerResult, 0, len(apiKeys))
	for _, apiKeyAgg := range apiKeys {
		p := apiKeyAgg.ToPrimitives()
		out = append(out, GetKeysByOwnerItem{
			ID:        p.ID,
			KeyName:   p.KeyName,
			Provider:  p.Provider,
			ModelName: p.ModelName,
			UsageEnv:  p.UsageEnv,
			Status:    p.Status,
			Balance:   p.Balance,
		})
	}

	return out, nil
}
