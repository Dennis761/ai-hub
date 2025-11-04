package models

import "time"

type APIKeyDoc struct {
	ID        string    `bson:"_id"`
	OwnerId   string    `bson:"ownerId"`
	KeyName   string    `bson:"keyName"`
	Provider  string    `bson:"provider"`
	ModelName string    `bson:"modelName"`
	KeyValue  string    `bson:"keyValue"`
	UsageEnv  string    `bson:"usageEnv"`
	Status    string    `bson:"status"`
	Balance   float64   `bson:"balance"`
	CreatedAt time.Time `bson:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt"`
}
