package projectrepo

import (
	"context"

	"ai_hub.com/app/core/domain/projectdomain"
	"ai_hub.com/app/infra/db/mongoose/mappers"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ProjectWriteRepoMongo struct {
	M *mongo.Collection
}

func NewProjectWriteRepoMongo(coll *mongo.Collection) *ProjectWriteRepoMongo {
	return &ProjectWriteRepoMongo{M: coll}
}

func (r *ProjectWriteRepoMongo) Create(ctx context.Context, entity *projectdomain.Project) (*projectdomain.Project, error) {
	doc, err := mappers.ProjectToPersistence(entity)
	if err != nil {
		return nil, err
	}

	if _, err := r.M.InsertOne(ctx, doc); err != nil {
		return nil, err
	}

	return mappers.ProjectFromDoc(doc)
}

func (r *ProjectWriteRepoMongo) Update(ctx context.Context, entity *projectdomain.Project) (*projectdomain.Project, error) {
	doc, err := mappers.ProjectToPersistence(entity)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": doc.ID}
	update := bson.M{
		"$set": bson.M{
			"name":        doc.Name,
			"status":      doc.Status,
			"apiKey":      doc.APIKey,
			"ownerId":     doc.OwnerID,
			"adminAccess": doc.AdminAccess,
			"updatedAt":   doc.UpdatedAt,
		},
	}

	opts := options.FindOneAndUpdate().
		SetReturnDocument(options.After).
		SetUpsert(false)

	var out mappers.ProjectDoc
	if err := r.M.FindOneAndUpdate(ctx, filter, update, opts).Decode(&out); err != nil {
		return nil, err
	}

	return mappers.ProjectFromDoc(&out)
}

func (r *ProjectWriteRepoMongo) Delete(ctx context.Context, id projectdomain.ProjectID) error {
	projectIDStr := id.Value()

	tasksColl := r.M.Database().Collection("tasks")
	promptsColl := r.M.Database().Collection("prompts")

	cursor, err := tasksColl.Find(ctx, bson.M{"projectId": projectIDStr}, options.Find().SetProjection(bson.M{"_id": 1}))
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	taskIDs := make([]string, 0, 8)
	for cursor.Next(ctx) {
		var row struct {
			ID string `bson:"_id"`
		}
		if err := cursor.Decode(&row); err != nil {
			return err
		}
		taskIDs = append(taskIDs, row.ID)
	}
	if err := cursor.Err(); err != nil {
		return err
	}

	// 2. delete the prompts of all these tasks
	if len(taskIDs) > 0 {
		if _, err := promptsColl.DeleteMany(ctx, bson.M{"taskId": bson.M{"$in": taskIDs}}); err != nil {
			return err
		}
	}

	// 3. delete the project tasks themselves
	if _, err := tasksColl.DeleteMany(ctx, bson.M{"projectId": projectIDStr}); err != nil {
		return err
	}

	// 4. delete the project itself
	_, err = r.M.DeleteOne(ctx, bson.M{"_id": projectIDStr})
	return err
}
