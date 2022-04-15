package token_repository

import (
	"context"
	"time"

	"github.com/nuntiodev/block-proto/go_block"
	"github.com/softcorp-io/x/cryptox"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	actionCreate = iota
)

const (
	expiresAfterIndex = "expires_after_index"
)

type Token struct {
	Id                      string    `bson:"_id" json:"_id"`
	UserId                  string    `bson:"user_id" json:"user_id"`
	Blocked                 bool      `bson:"blocked" json:"blocked"`
	Device                  string    `bson:"device" json:"device"`
	BlockedAt               time.Time `bson:"blocked_at" json:"blocked_at"`
	CreatedAt               time.Time `bson:"created_at" json:"created_at"`
	UsedAt                  time.Time `bson:"used_at" json:"used_at"`
	ExpiresAt               time.Time `bson:"expires_at" json:"expires_at"` // unix time
	Encrypted               bool      `bson:"encrypted" json:"encrypted"`
	InternalEncryptionLevel int       `bson:"internal_encryption_level" json:"internal_encryption_level"`
}

type TokenRepository interface {
	Create(ctx context.Context, token *go_block.Token) (*go_block.Token, error)
	Block(ctx context.Context, token *go_block.Token) (*go_block.Token, error)
	IsBlocked(ctx context.Context, token *go_block.Token) (bool, error)
	UpdateUsedAt(ctx context.Context, token *go_block.Token) (*go_block.Token, error)
	GetTokens(ctx context.Context, token *go_block.Token) ([]*go_block.Token, error)
	Get(ctx context.Context, token *go_block.Token) (*go_block.Token, error)
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
