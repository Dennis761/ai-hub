package apikeyapp

import (
	"context"

	"ai_hub.com/app/core/domain/apikeydomain"
	"ai_hub.com/app/core/ports/apikeyports"
)

type DeleteAPIKeyHandler struct {
	readRepo  apikeyports.APIKeyReadRepository
	writeRepo apikeyports.APIKeyWriteRepository
	uow       apikeyports.UnitOfWorkPort
}

func NewDeleteAPIKeyHandler(
	readRepo apikeyports.APIKeyReadRepository,
	writeRepo apikeyports.APIKeyWriteRepository,
	uow apikeyports.UnitOfWorkPort,
) *DeleteAPIKeyHandler {
	return &DeleteAPIKeyHandler{
		readRepo:  readRepo,
		writeRepo: writeRepo,
		uow:       uow,
	}
}

func (h *DeleteAPIKeyHandler) Delete(
	ctx context.Context,
	cmd DeleteAPIKeyCommand,
) error {
	_, err := h.uow.WithTransaction(ctx, func(txCtx context.Context) (any, error) {
		// normalize inputs
		keyID, err := apikeydomain.NewAPIKeyID(cmd.ID)
		if err != nil {
			return nil, err
		}
		ownerID, err := apikeydomain.NewOwnerID(cmd.OwnerID)
		if err != nil {
			return nil, err
		}

		// load aggregate
		apiKey, err := h.readRepo.FindByID(txCtx, keyID)
		if err != nil {
			return nil, err
		}
		if apiKey == nil {
			return nil, apikeydomain.APIKeyNotFound()
		}

		// ACL
		if apiKey.OwnerID().Value() != ownerID.Value() {
			return nil, apikeydomain.Forbidden()
		}

		// delete
		if err := h.writeRepo.Delete(txCtx, keyID); err != nil {
			return nil, err
		}

		return nil, nil
	})
	return err
}
