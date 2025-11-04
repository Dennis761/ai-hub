package apikeydomain

import "time"

type CreateAPIKeyProps struct {
	ID        APIKeyID
	OwnerID   OwnerID
	KeyName   APIKeyName
	KeyValue  EncryptedKeyValue
	Provider  ProviderName
	ModelName ModelName
	UsageEnv  UsageEnv
	Status    APIKeyStatus
	Balance   APIKeyBalance
	Now       *time.Time
}

// Create builds a new API key aggregate.
func Create(props CreateAPIKeyProps) (APIKey, error) {
	currentTime := time.Now().UTC()
	if props.Now != nil {
		currentTime = props.Now.UTC()
	}

	// default env
	usageEnv := "prod"
	if _, ok := allowedApiKeyEnviroments[usageEnv]; !ok {
		return APIKey{}, InvalidUsageEnv()
	}

	// default status
	apiKeyStatus := "active"
	if _, ok := allowedApiKeyStatuses[apiKeyStatus]; !ok {
		return APIKey{}, InvalidStatus()
	}

	return APIKey{
		id:          props.ID,
		ownerID:     props.OwnerID,
		name:        props.KeyName,
		keyValueEnc: props.KeyValue,
		provider:    props.Provider,
		model:       props.ModelName,
		usageEnv:    UsageEnv(usageEnv),
		status:      APIKeyStatus(apiKeyStatus),
		balance:     props.Balance,
		createdAt:   currentTime,
		updatedAt:   currentTime,
	}, nil
}

type RestoreProps struct {
	ID        APIKeyID
	OwnerID   OwnerID
	KeyName   APIKeyName
	KeyValue  EncryptedKeyValue
	Provider  ProviderName
	ModelName ModelName
	UsageEnv  string
	Status    string
	Balance   APIKeyBalance
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Restore rebuilds an API key aggregate from persisted data.
func Restore(props RestoreProps) (APIKey, error) {
	usageEnv := props.UsageEnv
	if _, ok := allowedApiKeyEnviroments[usageEnv]; !ok {
		return APIKey{}, InvalidUsageEnv()
	}

	apiKeyStatus := props.Status
	if _, ok := allowedApiKeyStatuses[apiKeyStatus]; !ok {
		return APIKey{}, InvalidStatus()
	}

	return APIKey{
		id:          props.ID,
		ownerID:     props.OwnerID,
		name:        props.KeyName,
		keyValueEnc: props.KeyValue,
		provider:    props.Provider,
		model:       props.ModelName,
		usageEnv:    UsageEnv(usageEnv),
		status:      APIKeyStatus(apiKeyStatus),
		balance:     props.Balance,
		createdAt:   props.CreatedAt,
		updatedAt:   props.UpdatedAt,
	}, nil
}
