package mappers

import (
	"time"

	"ai_hub.com/app/core/domain/apikeydomain"
)

type APIKeyDoc struct {
	ID        string    `bson:"_id"`
	OwnerID   string    `bson:"ownerId"`
	KeyName   string    `bson:"keyName"`
	KeyValue  string    `bson:"keyValue"`
	Provider  string    `bson:"provider"`
	ModelName string    `bson:"modelName"`
	UsageEnv  string    `bson:"usageEnv"`
	Status    string    `bson:"status"`
	Balance   float64   `bson:"balance"`
	CreatedAt time.Time `bson:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt"`
}

func APIKeyFromDoc(doc *APIKeyDoc) (*apikeydomain.APIKey, error) {
	if doc == nil {
		return nil, nil
	}

	// ID
	id, err := apikeydomain.NewAPIKeyID(doc.ID)
	if err != nil {
		return nil, err
	}

	// OwnerID
	owner, err := apikeydomain.NewOwnerID(doc.OwnerID)
	if err != nil {
		return nil, err
	}

	// KeyName
	name, err := apikeydomain.NewAPIKeyName(doc.KeyName)
	if err != nil {
		return nil, err
	}

	// KeyEnc
	keyEnc, err := apikeydomain.NewEncryptedKeyValueFromString(doc.KeyValue)
	if err != nil {
		return nil, err
	}

	// Provider
	provider, err := apikeydomain.NewProviderName(doc.Provider)
	if err != nil {
		return nil, err
	}

	// ModelName
	model, err := apikeydomain.NewModelName(doc.ModelName)
	if err != nil {
		return nil, err
	}

	// Balance
	balance := apikeydomain.NewAPIKeyBalance(doc.Balance)

	// UsageEnv
	usageEnv := apikeydomain.UsageEnv(doc.UsageEnv)

	// Status
	status := apikeydomain.APIKeyStatus(doc.Status)

	// CreatedAt & UpdatedAt
	createdAt := doc.CreatedAt
	updatedAt := doc.UpdatedAt

	props := apikeydomain.RestoreProps{
		ID:        id,
		OwnerID:   owner,
		KeyName:   name,
		KeyValue:  keyEnc,
		Provider:  provider,
		ModelName: model,
		UsageEnv:  string(usageEnv),
		Status:    string(status),
		Balance:   balance,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	ak, err := apikeydomain.Restore(props)
	if err != nil {
		return nil, err
	}
	return &ak, nil
}

func APIKeyToPersistence(entity *apikeydomain.APIKey) (*APIKeyDoc, error) {
	p := entity.ToPrimitives()

	createdAt, err := time.Parse(time.RFC3339, p.CreatedAt)
	if err != nil {
		return nil, err
	}
	updatedAt, err := time.Parse(time.RFC3339, p.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &APIKeyDoc{
		ID:        p.ID,
		OwnerID:   p.OwnerID,
		KeyName:   p.KeyName,
		KeyValue:  p.KeyValue,
		Provider:  p.Provider,
		ModelName: p.ModelName,
		UsageEnv:  p.UsageEnv,
		Status:    p.Status,
		Balance:   p.Balance,
		CreatedAt: createdAt.UTC(),
		UpdatedAt: updatedAt.UTC(),
	}, nil
}
