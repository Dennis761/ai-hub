package apikeyapp

type UpdateAPIKeyCommand struct {
	ID string `json:"_id"`

	OwnerID string `json:"ownerID"`

	KeyName   *string  `json:"keyName,omitempty"`
	KeyValue  *string  `json:"keyValue,omitempty"`
	Provider  *string  `json:"provider,omitempty"`
	ModelName *string  `json:"modelName,omitempty"`
	UsageEnv  *string  `json:"usageEnv,omitempty"`
	Status    *string  `json:"status,omitempty"`
	Balance   *float64 `json:"balance,omitempty"`
}
