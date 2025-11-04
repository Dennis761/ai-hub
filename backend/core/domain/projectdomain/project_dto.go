package projectdomain

import "time"

type ProjectPrimitives struct {
	ID          string   `json:"_id"`
	Name        string   `json:"name"`
	Status      string   `json:"status"`
	APIKey      string   `json:"apiKey"`
	OwnerID     string   `json:"ownerID"`
	AdminAccess []string `json:"adminAccess"`
	CreatedAt   string   `json:"createdAt"`
	UpdatedAt   string   `json:"updatedAt"`
}

func (p Project) ToPrimitives() ProjectPrimitives {
	return ProjectPrimitives{
		ID:          p.id.Value(),
		Name:        p.name.Value(),
		Status:      p.status,
		APIKey:      p.apiKey.Value(),
		OwnerID:     p.ownerID.Value(),
		AdminAccess: p.AdminAccess(),
		CreatedAt:   p.createdAt.UTC().Format(time.RFC3339),
		UpdatedAt:   p.updatedAt.UTC().Format(time.RFC3339),
	}
}
