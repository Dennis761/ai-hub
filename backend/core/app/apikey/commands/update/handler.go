package apikeyapp

import (
	"context"

	"ai_hub.com/app/core/domain/apikeydomain"
	"ai_hub.com/app/core/ports/apikeyports"
)

type UpdateAPIKeyHandler struct {
	readRepo  apikeyports.APIKeyReadRepository
	writeRepo apikeyports.APIKeyWriteRepository
	crypto    apikeyports.CryptoPort
	uow       apikeyports.UnitOfWorkPort
	catalog   apikeyports.ModelCatalogPort
}

func NewUpdateAPIKeyHandler(
	readRepo apikeyports.APIKeyReadRepository,
	writeRepo apikeyports.APIKeyWriteRepository,
	crypto apikeyports.CryptoPort,
	uow apikeyports.UnitOfWorkPort,
	catalog apikeyports.ModelCatalogPort,
) *UpdateAPIKeyHandler {
	return &UpdateAPIKeyHandler{
		readRepo:  readRepo,
		writeRepo: writeRepo,
		crypto:    crypto,
		uow:       uow,
		catalog:   catalog,
	}
}

func (h *UpdateAPIKeyHandler) Update(
	ctx context.Context,
	cmd UpdateAPIKeyCommand,
) (*apikeydomain.APIKey, error) {

	res, err := h.uow.WithTransaction(ctx, func(txCtx context.Context) (any, error) {
		// validate IDs
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

		// next name
		nextName := apiKey.KeyName()
		if cmd.KeyName != nil {
			nextName, err = apikeydomain.NewAPIKeyName(*cmd.KeyName)
			if err != nil {
				return nil, err
			}
		}

		// next env
		nextEnv := apiKey.UsageEnv()
		if cmd.UsageEnv != nil {
			nextEnv, err = apikeydomain.NewUsageEnv(*cmd.UsageEnv)
			if err != nil {
				return nil, err
			}
		}

		// next provider
		nextProvider := apiKey.Provider()
		if cmd.Provider != nil {
			nextProvider, err = apikeydomain.NewProviderName(*cmd.Provider)
			if err != nil {
				return nil, err
			}
		}

		// next model
		nextModel := apiKey.ModelName()
		if cmd.ModelName != nil {
			nextModel, err = apikeydomain.NewModelName(*cmd.ModelName)
			if err != nil {
				return nil, err
			}
		}

		// optional status
		var nextStatus apikeydomain.APIKeyStatus
		statusProvided := cmd.Status != nil
		if statusProvided {
			nextStatus, err = apikeydomain.NewAPIKeyStatus(*cmd.Status)
			if err != nil {
				return nil, err
			}
		}

		// optional balance
		var nextBalance apikeydomain.APIKeyBalance
		balanceProvided := cmd.Balance != nil
		if balanceProvided {
			nextBalance = apikeydomain.NewAPIKeyBalance(*cmd.Balance)
		}

		// optional new secret
		secretProvided := cmd.KeyValue != nil && *cmd.KeyValue != ""
		var incomingSecret string
		if secretProvided {
			plainVO, err := apikeydomain.NewPlainAPIKeyValue(*cmd.KeyValue)
			if err != nil {
				return nil, err
			}
			incomingSecret = plainVO.Value()
		}

		// model/provider changes → re-validate
		providerChanged := nextProvider.Value() != apiKey.Provider().Value()
		modelChanged := nextModel.Value() != apiKey.ModelName().Value()

		if providerChanged || modelChanged {
			secretForCheck := incomingSecret
			if !secretProvided {
				enc := apiKey.KeyValueEnc().ExposeForPersistence()
				secretForCheck, err = h.crypto.Decrypt(enc)
				if err != nil {
					return nil, err
				}
			}
			if _, err := h.catalog.GetModelInfoOrThrow(
				txCtx,
				nextProvider.Value(),
				nextModel.Value(),
				secretForCheck,
			); err != nil {
				return nil, err
			}
		} else if secretProvided {
			if _, err := h.catalog.GetModelInfoOrThrow(
				txCtx,
				nextProvider.Value(),
				nextModel.Value(),
				incomingSecret,
			); err != nil {
				return nil, err
			}
		}

		// uniqueness: name+env per owner
		nameChanged := nextName.Value() != apiKey.KeyName().Value()
		envChanged := nextEnv != apiKey.UsageEnv()

		if nameChanged || envChanged {
			sameNameKeys, err := h.readRepo.FindAllByKeyName(txCtx, nextName)
			if err != nil {
				return nil, err
			}

			// same name but different owner
			for _, k := range sameNameKeys {
				if k.OwnerID().Value() != ownerID.Value() {
					return nil, apikeydomain.KeyNameAlreadyUsedByAnotherUser(nextName.Value())
				}
			}

			// same owner + same name + same env → conflict
			for _, k := range sameNameKeys {
				if k.ID().Value() == apiKey.ID().Value() {
					continue
				}
				if k.OwnerID().Value() == ownerID.Value() &&
					k.KeyName().Value() == nextName.Value() &&
					k.UsageEnv() == nextEnv {
					return nil, apikeydomain.EnvAlreadyExistsForThisName(nextName.Value(), string(nextEnv))
				}
			}
		}

		changed := false

		if nameChanged {
			apiKey.Rename(nextName)
			changed = true
		}
		if providerChanged || modelChanged {
			apiKey.RebindProviderModel(nextProvider, nextModel)
			changed = true
		}
		if envChanged {
			if err := apiKey.MoveToEnvironment(nextEnv); err != nil {
				return nil, err
			}
			changed = true
		}
		if statusProvided {
			if err := apiKey.SetStatus(nextStatus); err != nil {
				return nil, err
			}
			changed = true
		}
		if balanceProvided {
			apiKey.SetBalance(nextBalance)
			changed = true
		}
		if secretProvided {
			enc, err := h.crypto.Encrypt(incomingSecret)
			if err != nil {
				return nil, err
			}
			encVO, err := apikeydomain.NewEncryptedKeyValueFromString(enc)
			if err != nil {
				return nil, err
			}
			apiKey.RotateEncryptedValue(encVO)
			changed = true
		}

		if !changed {
			return apiKey, nil
		}

		updated, err := h.writeRepo.Update(txCtx, apiKey)
		if err != nil {
			return nil, err
		}
		return updated, nil
	})
	if err != nil {
		return nil, err
	}

	return res.(*apikeydomain.APIKey), nil
}
