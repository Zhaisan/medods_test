package db

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const refreshTokensCollection = "refresh_tokens"
const dbName = "authDB"

func SaveRefreshToken(ctx context.Context, db *mongo.Client, token Token) error {
	collection := db.Database(dbName).Collection(refreshTokensCollection)
	_, err := collection.InsertOne(ctx, token)
	return err
}

func FindRefreshTokenByUserID(ctx context.Context, db *mongo.Client, userID string) (Token, error) {
	collection := db.Database(dbName).Collection(refreshTokensCollection)
	var token Token
	err := collection.FindOne(ctx, bson.M{"userId": userID}).Decode(&token)
	return token, err
}

func DeleteRefreshToken(ctx context.Context, db *mongo.Client, tokenHash string) error {
	collection := db.Database(dbName).Collection(refreshTokensCollection)
	_, err := collection.DeleteOne(ctx, bson.M{"tokenHash": tokenHash})
	return err
}