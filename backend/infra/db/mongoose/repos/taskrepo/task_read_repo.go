package taskrepo

import (
	"context"

	taskdomain "ai_hub.com/app/core/domain/taskdomain"
	taskports "ai_hub.com/app/core/ports/taskports"
	"ai_hub.com/app/infra/db/mongoose/mappers"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var _ taskports.TaskReadRepository = (*TaskReadRepoMongo)(nil)

type TaskReadRepoMongo struct {
	M *mongo.Collection
}

func NewTaskReadRepoMongo(coll *mongo.Collection) *TaskReadRepoMongo {
	return &TaskReadRepoMongo{M: coll}
}

func (r *TaskReadRepoMongo) FindByID(ctx context.Context, id taskdomain.TaskID) (*taskdomain.Task, error) {
	filter := bson.M{"_id": id.Value()}

	var doc mappers.TaskDoc
	err := r.M.FindOne(ctx, filter).Decode(&doc)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return mappers.TaskFromDoc(&doc)
}

func (r *TaskReadRepoMongo) FindAllByProject(ctx context.Context, projectId taskdomain.TaskProjectID) ([]taskdomain.Task, error) {
	filter := bson.M{"projectId": projectId.Value()}
	opts := options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}})

	cur, err := r.M.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var out []taskdomain.Task
	for cur.Next(ctx) {
		var d mappers.TaskDoc
		if err := cur.Decode(&d); err != nil {
			return nil, err
		}
		entPtr, err := mappers.TaskFromDoc(&d)
		if err != nil {
			return nil, err
		}
		if entPtr == nil {
			continue
		}
		out = append(out, *entPtr)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return out, nil
}
