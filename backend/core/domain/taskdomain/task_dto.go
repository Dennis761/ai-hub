package taskdomain

import "time"

type TaskPrimitives struct {
	ID          string  `json:"_id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
	ProjectID   string  `json:"projectID"`
	APIMethod   string  `json:"apiMethod"`
	Status      string  `json:"status"`
	CreatedBy   string  `json:"createdBy"`
	CreatedAt   string  `json:"createdAt"`
	UpdatedAt   string  `json:"updatedAt"`
}

func (t Task) ToPrimitives() TaskPrimitives {
	return TaskPrimitives{
		ID:          t.id.Value(),
		Name:        t.name.Value(),
		Description: t.description.Value(),
		ProjectID:   t.projectID.Value(),
		APIMethod:   t.apiMethod.Value(),
		Status:      t.status,
		CreatedBy:   t.createdBy.Value(),
		CreatedAt:   t.createdAt.UTC().Format(time.RFC3339),
		UpdatedAt:   t.updatedAt.UTC().Format(time.RFC3339),
	}
}
