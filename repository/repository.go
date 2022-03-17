package repository

import (
	"context"
	"github.com/softcorp-io/block-user-service/crypto"
	"github.com/softcorp-io/block-user-service/repository/user_repository"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type Repository interface {
	Liveness(ctx context.Context) error
	Users(ctx context.Context, namespace string) (user_repository.UserRepository, error)
}

type defaultRepository struct {
	namespace   string
	mongoClient *mongo.Client
	crypto      crypto.Crypto
}

func (r *defaultRepository) Liveness(ctx context.Context) error {
	if err := r.mongoClient.Ping(ctx, nil); err != nil {
		return err
	}
	return nil
}

func (r *defaultRepository) Users(ctx context.Context, namespace string) (user_repository.UserRepository, error) {
	if namespace == "" {
		namespace = "blocks-db"
	}
	collection := r.mongoClient.Database(namespace).Collection("users")
	userRepository, err := user_repository.NewUserRepository(ctx, collection, r.crypto)
	if err != nil {
		return nil, err
	}
	return userRepository, nil
}

func New(mongoClient *mongo.Client, crypto crypto.Crypto, zapLog *zap.Logger) (Repository, error) {
	zapLog.Info("creating repository_mock...")
	repository := &defaultRepository{
		mongoClient: mongoClient,
		crypto:      crypto,
	}
	return repository, nil
}
