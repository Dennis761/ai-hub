package promptapp

type RunPromptCommand struct {
	ID           string `json:"_id"`
	RunByAdminID string `json:"runByAdminID"`
}

type RunPromptResult struct {
	ID               string `json:"_id"`
	Version          int    `json:"version"`
	PromptText       string `json:"promptText"`
	ResponseText     string `json:"responseText"`
	ModelID          string `json:"modelID"`
	UpdatedAt        string `json:"updatedAt"`
	Cost             string `json:"cost"`             // fixed(6)
	BalanceRemaining string `json:"balanceRemaining"` // fixed(6)
}
