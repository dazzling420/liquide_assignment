package mongodb

import (
	"context"

	"liquide_assignment/internal/config"

	"go.mongodb.org/mongo-driver/mongo"
)

type MongoRepository struct {
	db     *mongo.Client
	ctx    context.Context
	config *config.Config
}

func InitMongoRepo(mongoClient *mongo.Client, ctx context.Context, config *config.Config) *MongoRepository {
	return &MongoRepository{
		db:     mongoClient,
		ctx:    ctx,
		config: config,
	}
}
