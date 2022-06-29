package config_repository

import (
	"context"
	"github.com/nuntiodev/hera-sdks/go_hera"
	"github.com/nuntiodev/hera/models"
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
	Create(ctx context.Context, config *go_hera.Config) error
	Get(ctx context.Context) (*models.Config, error)
	Update(ctx context.Context, config *go_hera.Config) error
	RegisterPublicKey(ctx context.Context, publicKey string) error
	RemovePublicKey(ctx context.Context) error
	Delete(ctx context.Context) error
}

type defaultConfigRepository struct {
	collection *mongo.Collection
	crypto     cryptox.Crypto
}

func New(ctx context.Context, collection *mongo.Collection, crypto cryptox.Crypto) (ConfigRepository, error) {
	return &defaultConfigRepository{
		collection: collection,
		crypto:     crypto,
	}, nil
}
