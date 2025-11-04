package promptapp

type UpdatePromptCommand struct {
	ID string `json:"id"`

	Name           *string `json:"name,omitempty"`
	ModelID        *string `json:"modelId,omitempty"`
	PromptText     *string `json:"promptText,omitempty"`
	ResponseText   *string `json:"responseText,omitempty"`
	ExecutionOrder *int    `json:"executionOrder,omitempty"`

	AdminID string `json:"adminId"`
}
