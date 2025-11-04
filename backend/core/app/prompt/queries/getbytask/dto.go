package promptapp

type GetPromptsByTaskQuery struct {
	TaskID  string `json:"taskID"`
	AdminID string `json:"adminID"`
}

type HistoryItem struct {
	Prompt    string `json:"prompt"`
	Response  string `json:"response"`
	Version   int    `json:"version"`
	CreatedAt string `json:"createdAt"`
}

type PromptListItem struct {
	ID             string        `json:"_id"`
	TaskID         string        `json:"taskId"`
	Name           string        `json:"name"`
	ModelID        string        `json:"modelId"`
	PromptText     string        `json:"promptText"`
	ResponseText   *string       `json:"responseText,omitempty"`
	History        []HistoryItem `json:"history,omitempty"`
	ExecutionOrder int           `json:"executionOrder"`
	Version        int           `json:"version"`
	CreatedAt      string        `json:"createdAt"`
	UpdatedAt      string        `json:"updatedAt"`
}
