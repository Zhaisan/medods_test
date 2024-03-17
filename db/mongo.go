package db

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
	"zhaisan-medods/utils"
)

func MongoDBClient(uri string) *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		utils.Logger.WithError(err).Fatal("Error to create MongoDB client")
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		utils.Logger.WithError(err).Fatal("Error to connect to MongoDB")

	}

	return client
}