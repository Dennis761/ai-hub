package apikeydomain

type UsageEnv string
type APIKeyStatus string

const (
	APIKeyStatusActive   APIKeyStatus = "active"
	APIKeyStatusInactive APIKeyStatus = "inactive"
)

var allowedApiKeyEnviroments = map[string]struct{}{
	"prod": {},
	"dev":  {},
	"test": {},
}

var allowedApiKeyStatuses = map[string]struct{}{
	"active":   {},
	"inactive": {},
	"archived": {},
}
