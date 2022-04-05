package repository

import (
	"context"
	"github.com/softcorp-io/block-user-service/repository/token_repository"
	"github.com/softcorp-io/x/cryptox"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type Repository interface {
	Liveness(ctx context.Context) error
	Users() UsersBuilder
	Tokens(ctx context.Context, namespace string) (token_repository.TokenRespository, error)
}

type defaultRepository struct {
	namespace      string
	mongoClient    *mongo.Client
	crypto         cryptox.Crypto
	encryptionKeys []string
}

func (r *defaultRepository) Liveness(ctx context.Context) error {
	if err := r.mongoClient.Ping(ctx, nil); err != nil {
		return err
	}
	return nil
}

func (r *defaultRepository) Tokens(ctx context.Context, namespace string) (token_repository.TokenRespository, error) {
	if namespace == "" {
		namespace = "blocks-db"
	}
	collection := r.mongoClient.Database(namespace).Collection("user_tokens")
	tokenRepository, err := token_repository.New(ctx, collection)
	if err != nil {
		return nil, err
	}
	return tokenRepository, nil
}

func New(mongoClient *mongo.Client, crypto cryptox.Crypto, encryptionKeys []string, zapLog *zap.Logger) (Repository, error) {
	zapLog.Info("creating repository...")
	repository := &defaultRepository{
		mongoClient:    mongoClient,
		crypto:         crypto,
		encryptionKeys: encryptionKeys,
	}
	return repository, nil
}
