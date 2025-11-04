package uow

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type MongoUnitOfWork struct {
	client *mongo.Client
}

func NewMongoUnitOfWork(client *mongo.Client) *MongoUnitOfWork {
	return &MongoUnitOfWork{client: client}
}

func (u *MongoUnitOfWork) WithTransaction(ctx context.Context, work func(txCtx context.Context) (any, error)) (any, error) {
	sess, err := u.client.StartSession()
	if err != nil {
		return nil, err
	}
	defer sess.EndSession(ctx)

	return sess.WithTransaction(ctx, func(sc mongo.SessionContext) (interface{}, error) {
		return work(sc)
	})
}
