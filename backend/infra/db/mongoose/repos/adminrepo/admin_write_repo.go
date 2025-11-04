package adminrepo

import (
	"context"

	admindomain "ai_hub.com/app/core/domain/admindomain"
	"ai_hub.com/app/infra/db/mongoose/mappers"
	"ai_hub.com/app/infra/db/mongoose/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AdminWriteRepoMongo struct {
	col *mongo.Collection
}

func NewAdminWriteRepoMongo(coll *mongo.Collection) *AdminWriteRepoMongo {
	return &AdminWriteRepoMongo{col: coll}
}

func (r *AdminWriteRepoMongo) Create(ctx context.Context, entity *admindomain.Admin) (*admindomain.Admin, error) {
	dto, err := mappers.AdminToPersistence(entity)
	if err != nil {
		return nil, err
	}
	if _, err := r.col.InsertOne(ctx, dto); err != nil {
		return nil, err
	}

	var out models.AdminDoc
	if err := r.col.FindOne(ctx, bson.M{"_id": dto.ID}).Decode(&out); err != nil {
		return nil, err
	}
	return mappers.AdminFromDoc(&out)
}

func (r *AdminWriteRepoMongo) Update(ctx context.Context, entity *admindomain.Admin) (*admindomain.Admin, error) {
	dto, err := mappers.AdminToPersistence(entity)
	if err != nil {
		return nil, err
	}

	update := bson.M{
		"$set": bson.M{
			"email":                   dto.Email,
			"name":                    dto.Name,
			"isVerified":              dto.IsVerified,
			"password":                dto.Password,
			"role":                    dto.Role,
			"verificationCode":        dto.VerificationCode,
			"verificationCodeExpires": dto.VerificationCodeExpires,
			"isResetCodeConfirmed":    dto.IsResetCodeConfirmed,
			"updatedAt":               dto.UpdatedAt,
		},
	}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After).SetUpsert(false)

	var out models.AdminDoc
	if err := r.col.FindOneAndUpdate(ctx, bson.M{"_id": dto.ID}, update, opts).Decode(&out); err != nil {
		return nil, err
	}
	return mappers.AdminFromDoc(&out)
}

func (r *AdminWriteRepoMongo) Delete(ctx context.Context, businessID string) error {
	_, err := r.col.DeleteOne(ctx, bson.M{"_id": businessID})
	return err
}

func (r *AdminWriteRepoMongo) Save(ctx context.Context, entity *admindomain.Admin) error {
	dto, err := mappers.AdminToPersistence(entity)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"email":                   dto.Email,
			"name":                    dto.Name,
			"isVerified":              dto.IsVerified,
			"password":                dto.Password,
			"role":                    dto.Role,
			"verificationCode":        dto.VerificationCode,
			"verificationCodeExpires": dto.VerificationCodeExpires,
			"isResetCodeConfirmed":    dto.IsResetCodeConfirmed,
			"updatedAt":               dto.UpdatedAt,
		},
		"$setOnInsert": bson.M{
			"createdAt": dto.CreatedAt,
		},
	}
	_, err = r.col.UpdateOne(ctx, bson.M{"_id": dto.ID}, update, options.Update().SetUpsert(true))
	return err
}
