package apikeyapp

import (
	"context"
	"time"

	"ai_hub.com/app/core/domain/apikeydomain"
	"ai_hub.com/app/core/ports/apikeyports"
)

type CreateAPIKeyHandler struct {
	readRepo     apikeyports.APIKeyReadRepository
	writeRepo    apikeyports.APIKeyWriteRepository
	crypto       apikeyports.CryptoPort
	idGen        apikeyports.IDGenerator
	uow          apikeyports.UnitOfWorkPort
	modelCatalog apikeyports.ModelCatalogPort
}

func NewCreateAPIKeyHandler(
	readRepo apikeyports.APIKeyReadRepository,
	writeRepo apikeyports.APIKeyWriteRepository,
	crypto apikeyports.CryptoPort,
	idGen apikeyports.IDGenerator,
	uow apikeyports.UnitOfWorkPort,
	modelCatalog apikeyports.ModelCatalogPort,
) *CreateAPIKeyHandler {
	return &CreateAPIKeyHandler{
		readRepo:     readRepo,
		writeRepo:    writeRepo,
		crypto:       crypto,
		idGen:        idGen,
		uow:          uow,
		modelCatalog: modelCatalog,
	}
}

func (h *CreateAPIKeyHandler) Create(
	ctx context.Context,
	cmd CreateAPIKeyCommand,
) (*apikeydomain.APIKey, error) {

	res, err := h.uow.WithTransaction(ctx, func(txCtx context.Context) (any, error) {
		// normalize inputs via VO
		owner, err := apikeydomain.NewOwnerID(cmd.OwnerID)
		if err != nil {
			return nil, err
		}
		provider, err := apikeydomain.NewProviderName(cmd.Provider)
		if err != nil {
			return nil, err
		}
		model, err := apikeydomain.NewModelName(cmd.ModelName)
		if err != nil {
			return nil, err
		}
		keyName, err := apikeydomain.NewAPIKeyName(cmd.KeyName)
		if err != nil {
			return nil, err
		}
		plainKeyVO, err := apikeydomain.NewPlainAPIKeyValue(cmd.KeyValue)
		if err != nil {
			return nil, err
		}
		plainKey := plainKeyVO.Value()

		envRaw := ""
		if cmd.UsageEnv != nil {
			envRaw = *cmd.UsageEnv
		}
		usageEnv, err := apikeydomain.NewUsageEnv(envRaw)
		if err != nil {
			return nil, err
		}

		statusRaw := ""
		if cmd.Status != nil {
			statusRaw = *cmd.Status
		}
		status, err := apikeydomain.NewAPIKeyStatus(statusRaw)
		if err != nil {
			return nil, err
		}

		// validate provider/model against catalog
		if _, err = h.modelCatalog.GetModelInfoOrThrow(
			txCtx,
			provider.Value(),
			model.Value(),
			plainKey,
		); err != nil {
			return nil, err
		}

		// uniqueness: name + env per owner
		existing, err := h.readRepo.FindAllByKeyName(txCtx, keyName)
		if err != nil {
			return nil, err
		}
		if err := ensureUniqueKeyNameAcrossOwners(
			existing,
			owner.Value(),
			keyName.Value(),
		); err != nil {
			return nil, err
		}
		if err := ensureUniqueEnvPerOwnerName(
			existing,
			owner.Value(),
			keyName.Value(),
			string(usageEnv),
		); err != nil {
			return nil, err
		}

		// encrypt key
		encrypted, err := h.crypto.Encrypt(plainKey)
		if err != nil {
			return nil, err
		}
		encVal, err := apikeydomain.NewEncryptedKeyValueFromString(encrypted)
		if err != nil {
			return nil, err
		}

		// id
		newID := h.idGen.NewID()
		idVO, err := apikeydomain.NewAPIKeyID(newID)
		if err != nil {
			return nil, err
		}

		// balance
		balance := apikeydomain.NewAPIKeyBalance(0)
		if cmd.Balance != nil {
			balance = apikeydomain.NewAPIKeyBalance(*cmd.Balance)
		}

		now := time.Now().UTC()

		// build aggregate
		apiKey, err := apikeydomain.Create(apikeydomain.CreateAPIKeyProps{
			ID:        idVO,
			OwnerID:   owner,
			KeyName:   keyName,
			KeyValue:  encVal,
			Provider:  provider,
			ModelName: model,
			UsageEnv:  usageEnv,
			Status:    status,
			Balance:   balance,
			Now:       &now,
		})
		if err != nil {
			return nil, err
		}

		// persist
		created, err := h.writeRepo.Create(txCtx, &apiKey)
		if err != nil {
			return nil, err
		}
		return created, nil
	})
	if err != nil {
		return nil, err
	}
	return res.(*apikeydomain.APIKey), nil
}

// ensure name not taken by another owner
func ensureUniqueKeyNameAcrossOwners(
	keys []*apikeydomain.APIKey,
	currentOwnerID string,
	keyName string,
) error {
	for _, k := range keys {
		if k.OwnerID().Value() != currentOwnerID {
			return apikeydomain.KeyNameAlreadyUsedByAnotherUser(keyName)
		}
	}
	return nil
}

// ensure env not duplicated for same owner+name
func ensureUniqueEnvPerOwnerName(
	keys []*apikeydomain.APIKey,
	currentOwnerID string,
	requestedKeyName string,
	requestedEnv string,
) error {
	for _, k := range keys {
		if k.OwnerID().Value() == currentOwnerID &&
			k.KeyName().Value() == requestedKeyName &&
			string(k.UsageEnv()) == requestedEnv {
			return apikeydomain.EnvAlreadyExistsForThisName(requestedKeyName, requestedEnv)
		}
	}
	return nil
}
