package token_repository

import (
	"context"
	"github.com/nuntiodev/hera-proto/go_hera"
	"github.com/nuntiodev/hera/models"
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
	Create(ctx context.Context, token *go_hera.Token) error
	Block(ctx context.Context, token *go_hera.Token) error
	IsBlocked(ctx context.Context, token *go_hera.Token) (bool, error)
	UpdateUsedAt(ctx context.Context, token *go_hera.Token) error
	GetTokens(ctx context.Context, token *go_hera.Token) ([]*models.Token, error)
	Get(ctx context.Context, token *go_hera.Token) (*models.Token, error)
	BuildIndexes(ctx context.Context) error
}

type mongodbRepository struct {
	collection *mongo.Collection
	crypto     cryptox.Crypto
}

func (r *mongodbRepository) BuildIndexes(ctx context.Context) error {
	expiresAtIndexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "expires_at", Value: 1},
		},
		Options: options.Index().SetExpireAfterSeconds(0).SetName(expiresAfterIndex),
	}
	if _, err := r.collection.Indexes().CreateOne(ctx, expiresAtIndexModel); err != nil {
		return err
	}
	return nil
}

func New(ctx context.Context, collection *mongo.Collection, crypto cryptox.Crypto) (TokenRepository, error) {
	return &mongodbRepository{
		collection: collection,
		crypto:     crypto,
	}, nil
}
