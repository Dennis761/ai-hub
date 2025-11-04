package apikeyapp

type GetKeyByIDQuery struct {
	ID      string `json:"_id"`
	OwnerID string `json:"ownerID"`
}

type GetKeyByIDResult struct {
	ID        string `json:"_id"`
	OwnerID   string `json:"ownerID"`
	KeyName   string `json:"keyName"`
	Provider  string `json:"provider"`
	ModelName string `json:"modelName"`
	UsageEnv  string `json:"usageEnv"`
	Status    string `json:"status"`
	KeyValue  string `json:"keyValue"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}
