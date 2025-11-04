package apikeyrepo

import (
	"context"

	"ai_hub.com/app/core/domain/apikeydomain"
	"ai_hub.com/app/infra/db/mongoose/mappers"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type APIKeyWriteRepoMongo struct {
	coll *mongo.Collection
}

func NewAPIKeyWriteRepoMongo(coll *mongo.Collection) *APIKeyWriteRepoMongo {
	return &APIKeyWriteRepoMongo{coll: coll}
}

func (r *APIKeyWriteRepoMongo) Create(ctx context.Context, entity *apikeydomain.APIKey) (*apikeydomain.APIKey, error) {
	doc, err := mappers.APIKeyToPersistence(entity)
	if err != nil {
		return nil, err
	}

	if _, err := r.coll.InsertOne(ctx, doc); err != nil {
		return nil, err
	}

	return mappers.APIKeyFromDoc(doc)
}

func (r *APIKeyWriteRepoMongo) Update(ctx context.Context, entity *apikeydomain.APIKey) (*apikeydomain.APIKey, error) {
	doc, err := mappers.APIKeyToPersistence(entity)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": doc.ID}
	update := bson.M{
		"$set": bson.M{
			"ownerId":   doc.OwnerID,
			"keyName":   doc.KeyName,
			"keyValue":  doc.KeyValue,
			"provider":  doc.Provider,
			"modelName": doc.ModelName,
			"usageEnv":  doc.UsageEnv,
			"status":    doc.Status,
			"balance":   doc.Balance,
			"updatedAt": doc.UpdatedAt,
		},
	}

	opts := options.FindOneAndUpdate().
		SetReturnDocument(options.After).
		SetUpsert(false)

	var out mappers.APIKeyDoc
	if err := r.coll.FindOneAndUpdate(ctx, filter, update, opts).Decode(&out); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return mappers.APIKeyFromDoc(&out)
}

func (r *APIKeyWriteRepoMongo) Delete(ctx context.Context, id apikeydomain.APIKeyID) error {
	_, err := r.coll.DeleteOne(ctx, bson.M{"_id": id.Value()})
	return err
}
