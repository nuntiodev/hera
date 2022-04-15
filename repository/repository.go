package repository

import (
	"context"
	"github.com/nuntio-dev/nuntio-user-block/repository/token_repository"
	"github.com/softcorp-io/x/cryptox"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type Repository interface {
	Liveness(ctx context.Context) error
	Users() UsersBuilder
	Tokens(ctx context.Context, namespace string) (token_repository.TokenRepository, error)
}

type defaultRepository struct {
	namespace              string
	mongodbClient          *mongo.Client
	crypto                 cryptox.Crypto
	internalEncryptionKeys []string
}

func (r *defaultRepository) Liveness(ctx context.Context) error {
	if err := r.mongodbClient.Ping(ctx, nil); err != nil {
		return err
	}
	return nil
}

func (r *defaultRepository) Tokens(ctx context.Context, namespace string) (token_repository.TokenRepository, error) {
	if namespace == "" {
		namespace = "blocks-db"
	}
	collection := r.mongodbClient.Database(namespace).Collection("user_tokens")
	tokenRepository, err := token_repository.New(ctx, collection, r.crypto, r.internalEncryptionKeys)
	if err != nil {
		return nil, err
	}
	return tokenRepository, nil
}

func New(mongoClient *mongo.Client, crypto cryptox.Crypto, encryptionKeys []string, zapLog *zap.Logger) (Repository, error) {
	zapLog.Info("creating repository...")
	repository := &defaultRepository{
		mongodbClient:          mongoClient,
		crypto:                 crypto,
		internalEncryptionKeys: encryptionKeys,
	}
	return repository, nil
}
