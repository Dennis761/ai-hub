package projectapp

type GetProjectByIDQuery struct {
	ProjectID string `json:"_id"`
	AdminID   string `json:"adminID"`
}

type GetProjectByIDResult struct {
	ID          string   `json:"_id"`
	Name        string   `json:"name"`
	Status      string   `json:"status"`
	OwnerID     string   `json:"ownerID"`
	AdminAccess []string `json:"adminAccess"`
	CreatedAt   string   `json:"createdAt"`
	UpdatedAt   string   `json:"updatedAt"`
}
