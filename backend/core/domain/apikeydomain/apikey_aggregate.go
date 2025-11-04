package apikeydomain

import "time"

// Api key — Aggregate Root.
type APIKey struct {
	id          APIKeyID
	ownerID     OwnerID
	name        APIKeyName
	keyValueEnc EncryptedKeyValue
	provider    ProviderName
	model       ModelName
	usageEnv    UsageEnv
	status      APIKeyStatus
	balance     APIKeyBalance
	createdAt   time.Time
	updatedAt   time.Time
}

func (a *APIKey) touch() {
	a.updatedAt = time.Now().UTC()
}
