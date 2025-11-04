package adminrepo

import (
	"context"
	"strings"

	admindomain "ai_hub.com/app/core/domain/admindomain"
	"ai_hub.com/app/infra/db/mongoose/mappers"
	"ai_hub.com/app/infra/db/mongoose/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AdminReadRepoMongo struct {
	C *mongo.Collection
}

func NewAdminReadRepoMongo(coll *mongo.Collection) *AdminReadRepoMongo {
	return &AdminReadRepoMongo{C: coll}
}

func (r *AdminReadRepoMongo) FindByID(ctx context.Context, id string) (*admindomain.Admin, error) {
	var doc models.AdminDoc
	err := r.C.FindOne(ctx, bson.M{"_id": id}).Decode(&doc)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return mappers.AdminFromDoc(&doc)
}

func (r *AdminReadRepoMongo) FindByEmail(ctx context.Context, email string) (*admindomain.Admin, error) {
	normalized := strings.ToLower(strings.TrimSpace(email))

	var doc models.AdminDoc
	err := r.C.FindOne(ctx, bson.M{"email": normalized}).Decode(&doc)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return mappers.AdminFromDoc(&doc)
}
