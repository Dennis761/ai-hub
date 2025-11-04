package models

import "time"

type TaskDoc struct {
	ID          string    `bson:"_id"`
	Name        string    `bson:"name"`
	Description *string   `bson:"description,omitempty"`
	ProjectId   string    `bson:"projectId"`
	APIMethod   string    `bson:"apiMethod"`
	Status      string    `bson:"status"`
	CreatedBy   string    `bson:"createdBy"`
	CreatedAt   time.Time `bson:"createdAt"`
	UpdatedAt   time.Time `bson:"updatedAt"`
}
