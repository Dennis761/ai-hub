// src/core/application/apikey/queries/get_by_id_handler.go
package apikeyapp

import (
	"context"

	"ai_hub.com/app/core/domain/apikeydomain"
	"ai_hub.com/app/core/ports/apikeyports"
)

type GetKeyByIDHandler struct {
	readRepo   apikeyports.APIKeyReadRepository
	cryptoPort apikeyports.CryptoPort
}

func NewGetKeyByIDHandler(
	readRepo apikeyports.APIKeyReadRepository,
	crypto apikeyports.CryptoPort,
) *GetKeyByIDHandler {
	return &GetKeyByIDHandler{
		readRepo:   readRepo,
		cryptoPort: crypto,
	}
}

func (h *GetKeyByIDHandler) GetKeyByID(
	ctx context.Context,
	q GetKeyByIDQuery,
) (*GetKeyByIDResult, error) {
	// normalize id
	idVO, err := apikeydomain.NewAPIKeyID(q.ID)
	if err != nil {
		return nil, err
	}

	// load aggregate
	apiKey, err := h.readRepo.FindByID(ctx, idVO)
	if err != nil {
		return nil, err
	}
	if apiKey == nil {
		return nil, apikeydomain.APIKeyNotFound()
	}

	// project safe dto
	p := apiKey.ToPrimitives()

	return &GetKeyByIDResult{
		ID:        p.ID,
		KeyName:   p.KeyName,
		Provider:  p.Provider,
		ModelName: p.ModelName,
		Status:    p.Status,
	}, nil
}
