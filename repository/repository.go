package repository

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"time"
)

const (
	defaultDb = "hera_db"
)

type Repository interface {
	Liveness(ctx context.Context) error
	UserRepositoryBuilder() UserRepositoryBuilder
	TokenRepositoryBuilder() TokenRepositoryBuilder
	ConfigRepositoryBuilder() ConfigRepositoryBuilder
	DropDatabase(ctx context.Context, namespace string) error
}

type defaultRepository struct {
	namespace               string
	mongodbClient           *mongo.Client
	internalEncryptionKeys  []string
	maxEmailVerificationAge time.Duration
}

func (r *defaultRepository) Liveness(ctx context.Context) error {
	if err := r.mongodbClient.Ping(ctx, nil); err != nil {
		return err
	}
	return nil
}

func (r *defaultRepository) DropDatabase(ctx context.Context, namespace string) error {
	if namespace == "" {
		return errors.New("missing required namespace")
	}
	return r.mongodbClient.Database(namespace).Drop(ctx)
}

func New(mongoClient *mongo.Client, encryptionKeys []string, logger *zap.Logger, maxEmailVerificationAge time.Duration) (Repository, error) {
	logger.Info("creating repository...")
	repository := &defaultRepository{
		mongodbClient:           mongoClient,
		internalEncryptionKeys:  encryptionKeys,
		maxEmailVerificationAge: maxEmailVerificationAge,
	}
	return repository, nil
}
