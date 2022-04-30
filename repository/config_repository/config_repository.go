package config_repository

import (
	"context"
	"github.com/nuntiodev/block-proto/go_block"
	"github.com/nuntiodev/x/cryptox"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	actionCreate = iota
	actionUpdate
)

type ConfigRepository interface {
	Create(ctx context.Context, config *go_block.Config) (*go_block.Config, error)
	Get(ctx context.Context, config *go_block.Config) (*go_block.Config, error)
	UpdateDetails(ctx context.Context, config *go_block.Config) (*go_block.Config, error)
	UpdateGeneralText(ctx context.Context, config *go_block.Config) (*go_block.Config, error)
	UpdateWelcomeText(ctx context.Context, config *go_block.Config) (*go_block.Config, error)
	UpdateRegisterText(ctx context.Context, config *go_block.Config) (*go_block.Config, error)
	UpdateLoginText(ctx context.Context, config *go_block.Config) (*go_block.Config, error)
	UpdateSettings(ctx context.Context, config *go_block.Config) (*go_block.Config, error)
	Delete(ctx context.Context, config *go_block.Config) error
}

type defaultConfigRepository struct {
	collection             *mongo.Collection
	crypto                 cryptox.Crypto
	internalEncryptionKeys []string
}

func newMongodbConfigRepository(ctx context.Context, collection *mongo.Collection, crypto cryptox.Crypto, internalEncryptionKeys []string) (*defaultConfigRepository, error) {
	return &defaultConfigRepository{
		collection:             collection,
		crypto:                 crypto,
		internalEncryptionKeys: internalEncryptionKeys,
	}, nil
}

func New(ctx context.Context, collection *mongo.Collection, crypto cryptox.Crypto, internalEncryptionKeys []string) (ConfigRepository, error) {
	return newMongodbConfigRepository(ctx, collection, crypto, internalEncryptionKeys)
}
