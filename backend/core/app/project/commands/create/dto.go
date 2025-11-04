package projectapp

type CreateProjectCommand struct {
	ID          *string  `json:"_id,omitempty"`
	Name        string   `json:"name"`
	APIKey      string   `json:"apiKey"`
	OwnerID     string   `json:"ownerID"`
	AdminAccess []string `json:"adminAccess"`
	Status      *string  `json:"status,omitempty"`
}
