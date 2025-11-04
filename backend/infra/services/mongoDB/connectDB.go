package mongo

import (
	"context"
	"log"
	"time"

	"ai_hub.com/app/infra/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func ConnectDB() *mongo.Client {
	if client != nil {
		return client
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(config.Env.MongoDBURI)

	c, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("[MongoDB] Connection error: %v", err)
	}

	if err := c.Ping(ctx, nil); err != nil {
		log.Fatalf("[MongoDB] Ping error: %v", err)
	}

	log.Println("[MongoDB] Connected successfully")

	client = c

	go func() {
		for {
			time.Sleep(5 * time.Second)
			if err := c.Ping(context.Background(), nil); err != nil {
				log.Println("[MongoDB] Disconnected or unreachable:", err)
			}
		}
	}()

	return client
}

func GetClient() *mongo.Client {
	if client == nil {
		return ConnectDB()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := client.Ping(ctx, nil); err != nil {
		log.Println("[MongoDB] Reconnecting due to ping failure:", err)
		client = nil
		return ConnectDB()
	}

	return client
}

func Disconnect(ctx context.Context) error {
	if client == nil {
		return nil
	}
	return client.Disconnect(ctx)
}
