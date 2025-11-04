package promptapp

type GetPromptByIDQuery struct {
	ID      string
	AdminID string
}

type HistoryItem struct {
	Prompt    string  `json:"prompt"`
	Response  *string `json:"response"`
	Version   int     `json:"version"`
	CreatedAt string  `json:"createdAt"`
}

type GetPromptByIDResult struct {
	ID             string        `json:"_id"`
	TaskID         string        `json:"taskID"`
	Name           string        `json:"name"`
	ModelID        string        `json:"modelID"`
	PromptText     string        `json:"promptText"`
	ResponseText   *string       `json:"responseText"`
	History        []HistoryItem `json:"history"`
	ExecutionOrder int           `json:"executionOrder"`
	Version        int           `json:"version"`
	CreatedAt      string        `json:"createdAt"`
	UpdatedAt      string        `json:"updatedAt"`
}
