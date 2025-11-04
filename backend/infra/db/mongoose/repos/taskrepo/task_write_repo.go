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

var _ taskports.TaskWriteRepository = (*TaskWriteRepoMongo)(nil)

type TaskWriteRepoMongo struct {
	M *mongo.Collection
}

func NewTaskWriteRepoMongo(coll *mongo.Collection) *TaskWriteRepoMongo {
	return &TaskWriteRepoMongo{M: coll}
}

func (r *TaskWriteRepoMongo) Create(ctx context.Context, entity taskdomain.Task) (*taskdomain.Task, error) {
	doc, err := mappers.TaskToPersistence(&entity)
	if err != nil {
		return nil, err
	}
	if _, err := r.M.InsertOne(ctx, doc); err != nil {
		return nil, err
	}
	return mappers.TaskFromDoc(doc)
}

func (r *TaskWriteRepoMongo) Update(ctx context.Context, entity taskdomain.Task) (*taskdomain.Task, error) {
	doc, err := mappers.TaskToPersistence(&entity)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": doc.ID}
	update := bson.M{
		"$set": bson.M{
			"name":        doc.Name,
			"description": doc.Description,
			"projectId":   doc.ProjectID,
			"apiMethod":   doc.APIMethod,
			"status":      doc.Status,
			"createdBy":   doc.CreatedBy,
			"updatedAt":   doc.UpdatedAt,
		},
	}

	opts := options.FindOneAndUpdate().
		SetReturnDocument(options.After).
		SetUpsert(false)

	var out mappers.TaskDoc
	if err := r.M.FindOneAndUpdate(ctx, filter, update, opts).Decode(&out); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return mappers.TaskFromDoc(&out)
}

func (r *TaskWriteRepoMongo) Delete(ctx context.Context, id taskdomain.TaskID) error {
	taskIDStr := id.Value()

	promptsColl := r.M.Database().Collection("prompts")

	// 1. delete all prompts of this task
	if _, err := promptsColl.DeleteMany(ctx, bson.M{"taskId": taskIDStr}); err != nil {
		return err
	}

	// 2. delete the task itself
	_, err := r.M.DeleteOne(ctx, bson.M{"_id": taskIDStr})
	return err
}
