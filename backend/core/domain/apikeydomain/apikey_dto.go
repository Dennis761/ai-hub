package apikeydomain

import "time"

type APIKeyPrimitives struct {
	ID        string  `json:"_id"`
	OwnerID   string  `json:"ownerID"`
	KeyName   string  `json:"keyName"`
	KeyValue  string  `json:"keyValue"`
	Provider  string  `json:"provider"`
	ModelName string  `json:"modelName"`
	UsageEnv  string  `json:"usageEnv"`
	Status    string  `json:"status"`
	Balance   float64 `json:"balance"`
	CreatedAt string  `json:"createdAt"`
	UpdatedAt string  `json:"updatedAt"`
}

func (a APIKey) ToPrimitives() APIKeyPrimitives {
	return APIKeyPrimitives{
		ID:        a.id.Value(),
		OwnerID:   a.ownerID.Value(),
		KeyName:   a.name.Value(),
		KeyValue:  a.keyValueEnc.ExposeForPersistence(),
		Provider:  a.provider.Value(),
		ModelName: a.model.Value(),
		UsageEnv:  string(a.usageEnv),
		Status:    string(a.status),
		Balance:   a.balance.ExposeForPersistence(),
		CreatedAt: a.createdAt.UTC().Format(time.RFC3339),
		UpdatedAt: a.updatedAt.UTC().Format(time.RFC3339),
	}
}
