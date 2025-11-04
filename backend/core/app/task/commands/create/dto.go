package taskapp

import "time"

type CreateTaskCommand struct {
	ID          *string
	Name        string
	Description *string
	ProjectID   string
	APIMethod   string
	Status      *string
	CreatedBy   string
	Now         *time.Time
}
