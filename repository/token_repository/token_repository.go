package token_repository

import (
	"context"
	"github.com/nuntiodev/block-proto/go_block"
	"github.com/nuntiodev/nuntio-user-block/models"
	"github.com/nuntiodev/x/cryptox"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	actionCreate = iota
	actionUpdate
)

const (
	expiresAfterIndex = "expires_after_index"
)

type TokenRepository interface {
	Create(ctx context.Context, token *go_block.Token) (*models.Token, error)
	Block(ctx context.Context, token *go_block.Token) (*models.Token, error)
	IsBlocked(ctx context.Context, token *go_block.Token) (bool, error)
	UpdateUsedAt(ctx context.Context, token *go_block.Token) (*models.Token, error)
	GetTokens(ctx context.Context, token *go_block.Token) ([]*models.Token, error)
	Get(ctx context.Context, token *go_block.Token) (*models.Token, error)
}

type mongodbRepository struct {
	collection             *mongo.Collection
	crypto                 cryptox.Crypto
	internalEncryptionKeys []string
}

func newMongodbTokenRepository(ctx context.Context, collection *mongo.Collection, crypto cryptox.Crypto, internalEncryptionKeys []string) (*mongodbRepository, error) {
	expiresAtIndexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "expires_at", Value: 1},
		},
		Options: options.Index().SetExpireAfterSeconds(0).SetName(expiresAfterIndex),
	}
	if _, err := collection.Indexes().CreateOne(ctx, expiresAtIndexModel); err != nil {
		return nil, err
	}
	return &mongodbRepository{
		collection:             collection,
		crypto:                 crypto,
		internalEncryptionKeys: internalEncryptionKeys,
	}, nil
}

func New(ctx context.Context, collection *mongo.Collection, crypto cryptox.Crypto, internalEncryptionKeys []string) (TokenRepository, error) {
	return newMongodbTokenRepository(ctx, collection, crypto, internalEncryptionKeys)
}
