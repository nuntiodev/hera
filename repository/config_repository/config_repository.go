package config_repository

import (
	"context"
	"github.com/nuntiodev/block-proto/go_block"
	"github.com/nuntiodev/nuntio-user-block/models"
	"github.com/nuntiodev/x/cryptox"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	actionCreate = iota
	actionUpdate
)

const (
	namespaceConfigName = "namespace_default_config"
)

type ConfigRepository interface {
	Create(ctx context.Context, config *go_block.Config) (*models.Config, error)
	GetNamespaceConfig(ctx context.Context) (*models.Config, error)
	UpdateDetails(ctx context.Context, config *go_block.Config) (*models.Config, error)
	UpdateSettings(ctx context.Context, config *go_block.Config) (*models.Config, error)
	Delete(ctx context.Context) error
}

type defaultConfigRepository struct {
	collection *mongo.Collection
	crypto     cryptox.Crypto
}

func newMongodbConfigRepository(ctx context.Context, collection *mongo.Collection, crypto cryptox.Crypto) (*defaultConfigRepository, error) {
	return &defaultConfigRepository{
		collection: collection,
		crypto:     crypto,
	}, nil
}

func New(ctx context.Context, collection *mongo.Collection, crypto cryptox.Crypto) (ConfigRepository, error) {
	return newMongodbConfigRepository(ctx, collection, crypto)
}
