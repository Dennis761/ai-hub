package promptrepo

import (
	"context"

	"ai_hub.com/app/core/domain/promptdomain"
	"ai_hub.com/app/infra/db/mongoose/mappers"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PromptWriteRepoMongo struct {
	M *mongo.Collection
}

func NewPromptWriteRepoMongo(coll *mongo.Collection) *PromptWriteRepoMongo {
	return &PromptWriteRepoMongo{M: coll}
}

func (r *PromptWriteRepoMongo) Create(ctx context.Context, entity *promptdomain.Prompt) (*promptdomain.Prompt, error) {
	doc, err := mappers.PromptToPersistence(entity)
	if err != nil {
		return nil, err
	}
	if _, err := r.M.InsertOne(ctx, doc); err != nil {
		return nil, err
	}
	return mappers.PromptFromDoc(doc)
}

func (r *PromptWriteRepoMongo) Update(ctx context.Context, entity *promptdomain.Prompt) (*promptdomain.Prompt, error) {
	doc, err := mappers.PromptToPersistence(entity)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": doc.ID}
	update := bson.M{
		"$set": bson.M{
			"name":           doc.Name,
			"modelId":        doc.ModelID,
			"promptText":     doc.PromptText,
			"responseText":   doc.ResponseText,
			"history":        doc.History,
			"executionOrder": doc.ExecutionOrder,
			"version":        doc.Version,
			"updatedAt":      doc.UpdatedAt,
		},
	}

	opts := options.FindOneAndUpdate().
		SetReturnDocument(options.After).
		SetUpsert(false)

	var out mappers.PromptDoc
	if err := r.M.FindOneAndUpdate(ctx, filter, update, opts).Decode(&out); err != nil {
		return nil, err
	}
	return mappers.PromptFromDoc(&out)
}

func (r *PromptWriteRepoMongo) UpdateMany(ctx context.Context, prompts []*promptdomain.Prompt) error {
	if len(prompts) == 0 {
		return nil
	}

	models := make([]mongo.WriteModel, 0, len(prompts))

	for _, p := range prompts {
		doc, err := mappers.PromptToPersistence(p)
		if err != nil {
			return err
		}
		m := mongo.NewUpdateOneModel().
			SetFilter(bson.M{"_id": doc.ID}).
			SetUpdate(bson.M{
				"$set": bson.M{
					"executionOrder": doc.ExecutionOrder,
					"updatedAt":      doc.UpdatedAt,
				},
			}).
			SetUpsert(false)
		models = append(models, m)
	}

	_, err := r.M.BulkWrite(ctx, models, options.BulkWrite())
	return err
}

func (r *PromptWriteRepoMongo) Delete(ctx context.Context, id promptdomain.PromptID) error {
	_, err := r.M.DeleteOne(ctx, bson.M{"_id": id.Value()})
	return err
}
