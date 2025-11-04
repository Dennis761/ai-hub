package apikeyrepo

import (
	"context"

	"ai_hub.com/app/core/domain/apikeydomain"
	"ai_hub.com/app/infra/db/mongoose/mappers"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type APIKeyReadRepoMongo struct {
	coll *mongo.Collection
}

func NewAPIKeyReadRepoMongo(coll *mongo.Collection) *APIKeyReadRepoMongo {
	return &APIKeyReadRepoMongo{coll: coll}
}

func (r *APIKeyReadRepoMongo) FindAllByOwner(ctx context.Context, ownerId apikeydomain.OwnerID) ([]*apikeydomain.APIKey, error) {
	opts := options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}})
	cur, err := r.coll.Find(ctx, bson.M{"ownerId": ownerId.Value()}, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var out []*apikeydomain.APIKey
	for cur.Next(ctx) {
		var doc mappers.APIKeyDoc
		if err := cur.Decode(&doc); err != nil {
			return nil, err
		}
		entity, err := mappers.APIKeyFromDoc(&doc)
		if err != nil {
			return nil, err
		}
		out = append(out, entity)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *APIKeyReadRepoMongo) FindByID(ctx context.Context, id apikeydomain.APIKeyID) (*apikeydomain.APIKey, error) {
	var doc mappers.APIKeyDoc
	if err := r.coll.FindOne(ctx, bson.M{"_id": id.Value()}).Decode(&doc); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return mappers.APIKeyFromDoc(&doc)
}

func (r *APIKeyReadRepoMongo) FindAllByKeyName(ctx context.Context, keyName apikeydomain.APIKeyName) ([]*apikeydomain.APIKey, error) {
	cur, err := r.coll.Find(ctx, bson.M{"keyName": keyName.Value()})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var out []*apikeydomain.APIKey
	for cur.Next(ctx) {
		var doc mappers.APIKeyDoc
		if err := cur.Decode(&doc); err != nil {
			return nil, err
		}
		entity, err := mappers.APIKeyFromDoc(&doc)
		if err != nil {
			return nil, err
		}
		out = append(out, entity)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *APIKeyReadRepoMongo) GetModelByID(ctx context.Context, modelID string) (*apikeydomain.APIKey, error) {
	var doc mappers.APIKeyDoc
	if err := r.coll.FindOne(ctx, bson.M{"_id": modelID}).Decode(&doc); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return mappers.APIKeyFromDoc(&doc)
}
