package database

import (
	"context"
	"liquide_assignment/internal/config"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func MongoConnect(mongoConfig config.MongoDb) (*mongo.Client, context.Context, error) {
	// Use a timeout context only for the connection attempt
	connectCtx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	//uri := fmt.Sprintf("mongodb://%s:%s@%s", mongoConfig.User, mongoConfig.Password, mongoConfig.Host)
	uri := "mongodb://localhost:27017"
	client, err := mongo.Connect(connectCtx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, nil, err
	}

	// Return a background context for ongoing operations
	// This context won't be cancelled and can be used for database operations
	return client, context.Background(), nil
}
