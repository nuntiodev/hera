package token_repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	expiresAfterIndex = "expires_after_index"
)

type Token struct {
	Id        string `bson:"_id" json:"id"`
	ExpiresAt int64  `bson:"expires_at" json:"expires_at"` // unix time
}

type TokenRespository interface {
	BlockToken(ctx context.Context, token *Token) error
	IsBlocked(ctx context.Context, token *Token) error
}

type mongoRepository struct {
	collection *mongo.Collection
}

func New(ctx context.Context, collection *mongo.Collection) (TokenRespository, error) {
	emailNamespaceIndexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "expires_at", Value: 1},
		},
		Options: options.Index().SetExpireAfterSeconds(0).SetName(expiresAfterIndex),
	}
	if _, err := collection.Indexes().CreateOne(ctx, emailNamespaceIndexModel); err != nil {
		return nil, err
	}
	return &mongoRepository{
		collection: collection,
	}, nil
}
