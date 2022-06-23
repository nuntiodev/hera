package repository

import (
	"context"
	"github.com/nuntiodev/hera/repository/config_repository"
	"github.com/nuntiodev/x/cryptox"
	"go.mongodb.org/mongo-driver/mongo"
)

type ConfigRepositoryBuilder interface {
	SetNamespace(namespace string) ConfigRepositoryBuilder
	Build(ctx context.Context) (config_repository.ConfigRepository, error)
}

type configRepositoryBuilder struct {
	internalEncryptionKeys []string
	namespace              string
	client                 *mongo.Client
}

func (cb *configRepositoryBuilder) SetNamespace(namespace string) ConfigRepositoryBuilder {
	cb.namespace = namespace
	return cb
}

func (cb *configRepositoryBuilder) Build(ctx context.Context) (config_repository.ConfigRepository, error) {
	if cb.namespace == "" {
		cb.namespace = defaultDb
	}
	crypto, err := cryptox.New(cb.internalEncryptionKeys, []string{})
	if err != nil {
		return nil, err
	}
	collection := cb.client.Database(cb.namespace).Collection("hera_config")
	configRepository, err := config_repository.New(ctx, collection, crypto)
	if err != nil {
		return nil, err
	}
	return configRepository, nil
}

func (r *defaultRepository) ConfigRepositoryBuilder() ConfigRepositoryBuilder {
	return &configRepositoryBuilder{
		client:                 r.mongodbClient,
		internalEncryptionKeys: r.internalEncryptionKeys,
	}
}
