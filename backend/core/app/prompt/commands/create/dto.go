package promptapp

type CreatePromptCommand struct {
	ID             *string `json:"_id"`
	TaskID         string  `json:"taskId"`
	Name           string  `json:"name"`
	ModelID        string  `json:"modelId"`
	PromptText     string  `json:"promptText"`
	ResponseText   *string `json:"responseText"`
	ExecutionOrder *int    `json:"executionOrder"`
	CreatedBy      string  `json:"createdBy"`
}
