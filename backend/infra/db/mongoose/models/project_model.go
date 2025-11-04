package models

import "time"

type ProjectDoc struct {
	ID          string    `bson:"_id"`
	Name        string    `bson:"name"`
	Status      string    `bson:"status"`
	APIKey      string    `bson:"apiKey"`
	OwnerId     string    `bson:"ownerId"`
	AdminAccess []string  `bson:"adminAccess"`
	CreatedAt   time.Time `bson:"createdAt"`
	UpdatedAt   time.Time `bson:"updatedAt"`
}
