package database

import (
	"context"
	"fmt"
	"log"

	"crud-users/internal/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Connect(ctx context.Context, cfg config.Config) (*mongo.Client, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		return nil, fmt.Errorf("connect to mongo: %w", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("ping mongo: %w", err)
	}

	log.Println("connected to mongo")
	return client, nil
}

func Close(ctx context.Context, client *mongo.Client) {
	if client == nil {
		return
	}
	if err := client.Disconnect(ctx); err != nil {
		log.Printf("failed to disconnect mongo: %v", err)
	}
}
