package projectapp

type UpdateProjectCommand struct {
	ID      string  `json:"_id"`
	OwnerID string  `json:"ownerID"`
	Name    *string `json:"name,omitempty"`
	APIKey  *string `json:"apiKey,omitempty"`
	Status  *string `json:"status,omitempty"`
}
