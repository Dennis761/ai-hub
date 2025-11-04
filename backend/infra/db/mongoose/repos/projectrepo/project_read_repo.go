package projectrepo

import (
	"context"

	"ai_hub.com/app/core/domain/projectdomain"
	"ai_hub.com/app/infra/db/mongoose/mappers"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ProjectReadRepoMongo struct {
	M *mongo.Collection
}

func NewProjectReadRepoMongo(coll *mongo.Collection) *ProjectReadRepoMongo {
	return &ProjectReadRepoMongo{M: coll}
}

func (r *ProjectReadRepoMongo) FindByID(ctx context.Context, id projectdomain.ProjectID) (*projectdomain.Project, error) {
	filter := bson.M{"_id": id.Value()}

	var doc mappers.ProjectDoc
	err := r.M.FindOne(ctx, filter).Decode(&doc)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return mappers.ProjectFromDoc(&doc)
}

func (r *ProjectReadRepoMongo) FindByName(ctx context.Context, name projectdomain.ProjectName) (*projectdomain.Project, error) {
	filter := bson.M{"name": name.Value()}

	var doc mappers.ProjectDoc
	err := r.M.FindOne(ctx, filter).Decode(&doc)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return mappers.ProjectFromDoc(&doc)
}

func (r *ProjectReadRepoMongo) FindAllAccessibleByAdmin(ctx context.Context, adminId string) ([]*projectdomain.Project, error) {
	filter := bson.M{
		"$or": []bson.M{
			{"ownerId": adminId},
			{"adminAccess": adminId},
		},
	}
	opts := options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}})

	cur, err := r.M.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var out []*projectdomain.Project
	for cur.Next(ctx) {
		var doc mappers.ProjectDoc
		if err := cur.Decode(&doc); err != nil {
			return nil, err
		}
		ent, err := mappers.ProjectFromDoc(&doc)
		if err != nil {
			return nil, err
		}
		out = append(out, ent)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return out, nil
}
