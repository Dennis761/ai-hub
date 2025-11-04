package apikeyapp

type CreateAPIKeyCommand struct {
	ID        *string  `json:"_id,omitempty"`
	OwnerID   string   `json:"ownerID"`
	KeyName   string   `json:"keyName"`
	KeyValue  string   `json:"keyValue"`
	Provider  string   `json:"provider"`
	ModelName string   `json:"modelName"`
	UsageEnv  *string  `json:"usageEnv,omitempty"`
	Status    *string  `json:"status,omitempty"`
	Balance   *float64 `json:"balance,omitempty"`
}
