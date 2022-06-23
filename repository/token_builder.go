package repository

import (
	"context"
	"github.com/nuntiodev/hera/repository/token_repository"
	"github.com/nuntiodev/x/cryptox"
	"go.mongodb.org/mongo-driver/mongo"
)

type TokenRepositoryBuilder interface {
	SetNamespace(namespace string) TokenRepositoryBuilder
	Build(ctx context.Context) (token_repository.TokenRepository, error)
}

type tokenRepositoryBuilder struct {
	namespace              string
	internalEncryptionKeys []string
	client                 *mongo.Client
}

func (tb *tokenRepositoryBuilder) SetNamespace(namespace string) TokenRepositoryBuilder {
	tb.namespace = namespace
	return tb
}

func (ub *tokenRepositoryBuilder) Build(ctx context.Context) (token_repository.TokenRepository, error) {
	if ub.namespace == "" {
		ub.namespace = defaultDb
	}
	crypto, err := cryptox.New(ub.internalEncryptionKeys, []string{})
	if err != nil {
		return nil, err
	}
	collection := ub.client.Database(ub.namespace).Collection("hera_tokens")
	tokenRepository, err := token_repository.New(ctx, collection, crypto)
	if err != nil {
		return nil, err
	}
	return tokenRepository, nil
}

func (r *defaultRepository) TokenRepositoryBuilder() TokenRepositoryBuilder {
	return &tokenRepositoryBuilder{
		client:                 r.mongodbClient,
		internalEncryptionKeys: r.internalEncryptionKeys,
	}
}
