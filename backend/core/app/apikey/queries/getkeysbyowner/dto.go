package apikeyapp

type GetKeysByOwnerQuery struct {
	OwnerID string
}

type GetKeysByOwnerItem struct {
	ID        string  `json:"id"`
	OwnerID   string  `json:"ownerId"`
	KeyName   string  `json:"keyName"`
	Provider  string  `json:"provider"`
	ModelName string  `json:"modelName"`
	UsageEnv  string  `json:"usageEnv"`
	Status    string  `json:"status"`
	Balance   float64 `json:"balance"`
	CreatedAt string  `json:"createdAt"`
	UpdatedAt string  `json:"updatedAt"`
}

type GetKeysByOwnerResult []GetKeysByOwnerItem
