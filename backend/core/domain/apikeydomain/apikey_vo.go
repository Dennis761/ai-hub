package apikeydomain

import (
	"math"
	"strings"
)

const (
	MinApiKeyNameLen = 6
	MaxApiKeyNameLen = 100
)

// ===================== APIKeyID =====================

type APIKeyID struct {
	value string
}

func NewAPIKeyID(raw string) (APIKeyID, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return APIKeyID{}, APIKeyIDRequired()
	}
	return APIKeyID{value: trimmed}, nil
}

func (id APIKeyID) Value() string { return id.value }

// ===================== OwnerID =====================

type OwnerID struct {
	value string
}

func NewOwnerID(raw string) (OwnerID, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return OwnerID{}, OwnerIDRequired()
	}
	return OwnerID{value: trimmed}, nil
}

func (id OwnerID) Value() string { return id.value }

// ===================== APIKeyName =====================

type APIKeyName struct {
	value string
}

func NewAPIKeyName(raw string) (APIKeyName, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" || len(trimmed) < MinApiKeyNameLen {
		return APIKeyName{}, InvalidAPIKeyName("")
	}
	if len(trimmed) > MaxApiKeyNameLen {
		return APIKeyName{}, APIKeyNameTooLong(100)
	}
	return APIKeyName{value: trimmed}, nil
}

func (n APIKeyName) Value() string { return n.value }

// ===================== PlainAPIKeyValue (NEW) =====================

// PlainAPIKeyValue — raw (unencrypted) key received from the client.
type PlainAPIKeyValue struct {
	value string
}

func NewPlainAPIKeyValue(raw string) (PlainAPIKeyValue, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return PlainAPIKeyValue{}, InvalidPlainAPIKeyValue()
	}
	return PlainAPIKeyValue{value: trimmed}, nil
}

func (p PlainAPIKeyValue) Value() string {
	return p.value
}

// ===================== EncryptedKeyValue =====================

type EncryptedKeyValue struct {
	blob string
}

func NewEncryptedKeyValueFromString(encrypted string) (EncryptedKeyValue, error) {
	if encrypted == "" || !strings.Contains(encrypted, ":") {
		return EncryptedKeyValue{}, InvalidEncryptedKeyFormat()
	}
	return EncryptedKeyValue{blob: encrypted}, nil
}

func (e EncryptedKeyValue) ExposeForPersistence() string { return e.blob }

// ===================== ProviderName =====================

type ProviderName struct {
	value string
}

func NewProviderName(raw string) (ProviderName, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return ProviderName{}, ProviderRequired()
	}
	if len(trimmed) > 100 {
		return ProviderName{}, ProviderTooLong(100)
	}
	return ProviderName{value: trimmed}, nil
}

func (p ProviderName) Value() string { return p.value }

// ===================== ModelName =====================

type ModelName struct {
	value string
}

func NewModelName(raw string) (ModelName, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return ModelName{}, ModelNameRequired()
	}
	if len(trimmed) > 120 {
		return ModelName{}, ModelNameTooLong(120)
	}
	return ModelName{value: trimmed}, nil
}

func (m ModelName) Value() string { return m.value }

// ===================== UsageEnv =====================

func NewUsageEnv(raw string) (UsageEnv, error) {
	s := strings.ToLower(strings.TrimSpace(raw))

	if _, ok := allowedApiKeyEnviroments[s]; !ok {
		return "", InvalidUsageEnv()
	}

	return UsageEnv(s), nil
}

// ===================== APIKeyStatus =====================

func NewAPIKeyStatus(raw string) (APIKeyStatus, error) {
	s := strings.ToLower(strings.TrimSpace(raw))

	if _, ok := allowedApiKeyStatuses[s]; !ok {
		return "", InvalidStatus()
	}

	return APIKeyStatus(s), nil
}

// ===================== APIKeyBalance =====================

type APIKeyBalance struct {
	value float64
}

func NewAPIKeyBalance(raw float64) APIKeyBalance {
	val := raw
	if math.IsNaN(val) || math.IsInf(val, 0) || val < 0 {
		val = 0
	}
	return APIKeyBalance{value: val}
}

func (b APIKeyBalance) ExposeForPersistence() float64 { return b.value }
func (b APIKeyBalance) Value() float64                { return b.value }
