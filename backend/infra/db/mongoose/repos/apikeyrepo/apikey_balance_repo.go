package apikeyrepo

import (
	"context"
	"time"

	apikeyports "ai_hub.com/app/core/ports/apikeyports"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var _ apikeyports.APIKeyBalancePort = (*APIKeyBalanceRepoMongo)(nil)

type APIKeyBalanceRepoMongo struct {
	coll *mongo.Collection
}

func NewAPIKeyBalanceRepoMongo(coll *mongo.Collection) *APIKeyBalanceRepoMongo {
	return &APIKeyBalanceRepoMongo{coll: coll}
}

func (r *APIKeyBalanceRepoMongo) SetBalance(
	ctx context.Context,
	apiKeyID string,
	newBalance float64,
	adminID string,
) error {
	update := bson.M{
		"$set": bson.M{
			"balance":   newBalance,
			"updatedAt": time.Now().UTC(),
			"updatedBy": adminID,
		},
	}

	_, err := r.coll.UpdateOne(ctx, bson.M{"_id": apiKeyID}, update)
	return err
}
