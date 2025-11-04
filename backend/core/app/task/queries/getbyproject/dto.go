package taskapp

type GetTasksByProjectQuery struct {
	ProjectID string `json:"projectID"`
	AdminID   string `json:"adminID"`
}

type TaskListItem struct {
	ID          string   `json:"_id"`
	Name        string   `json:"name"`
	Description *string  `json:"description,omitempty"`
	ProjectID   string   `json:"projectID"`
	APIMethod   string   `json:"apiMethod"`
	Status      string   `json:"status"`
	Prompts     []string `json:"prompts"`
	CreatedBy   string   `json:"createdBy"`
	CreatedAt   string   `json:"createdAt"`
	UpdatedAt   string   `json:"updatedAt"`
}
