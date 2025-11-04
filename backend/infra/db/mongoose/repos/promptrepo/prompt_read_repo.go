package promptrepo

import (
	"context"

	"ai_hub.com/app/core/domain/promptdomain"
	"ai_hub.com/app/core/domain/taskdomain"
	"ai_hub.com/app/infra/db/mongoose/mappers"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PromptReadRepoMongo struct {
	M *mongo.Collection
}

func NewPromptReadRepoMongo(coll *mongo.Collection) *PromptReadRepoMongo {
	return &PromptReadRepoMongo{M: coll}
}

func (r *PromptReadRepoMongo) FindByID(ctx context.Context, id promptdomain.PromptID) (*promptdomain.Prompt, error) {
	filter := bson.M{"_id": id.Value()}
	var doc mappers.PromptDoc
	err := r.M.FindOne(ctx, filter).Decode(&doc)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return mappers.PromptFromDoc(&doc)
}

func (r *PromptReadRepoMongo) FindAllByTask(ctx context.Context, taskID taskdomain.TaskID) ([]*promptdomain.Prompt, error) {
	filter := bson.M{"taskId": taskID.Value()}
	opts := options.Find().SetSort(bson.D{
		{Key: "executionOrder", Value: 1},
		{Key: "createdAt", Value: 1},
	})

	cur, err := r.M.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var out []*promptdomain.Prompt
	for cur.Next(ctx) {
		var d mappers.PromptDoc
		if err := cur.Decode(&d); err != nil {
			return nil, err
		}
		ent, err := mappers.PromptFromDoc(&d)
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

func (r *PromptReadRepoMongo) FindByIDs(ctx context.Context, ids []promptdomain.PromptID) ([]*promptdomain.Prompt, error) {
	if len(ids) == 0 {
		return []*promptdomain.Prompt{}, nil
	}

	raw := make([]string, 0, len(ids))
	for _, id := range ids {
		raw = append(raw, id.Value())
	}

	filter := bson.M{"_id": bson.M{"$in": raw}}
	cur, err := r.M.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var out []*promptdomain.Prompt
	for cur.Next(ctx) {
		var d mappers.PromptDoc
		if err := cur.Decode(&d); err != nil {
			return nil, err
		}
		ent, err := mappers.PromptFromDoc(&d)
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
