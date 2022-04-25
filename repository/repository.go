package repository

import (
	"context"
	"github.com/nuntiodev/nuntio-user-block/repository/config_repository"
	"github.com/nuntiodev/nuntio-user-block/repository/measurement_repository"
	"github.com/nuntiodev/nuntio-user-block/repository/token_repository"
	"github.com/nuntiodev/x/cryptox"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type Repository interface {
	Liveness(ctx context.Context) error
	Users() UsersBuilder
	Tokens(ctx context.Context, namespace string) (token_repository.TokenRepository, error)
	Measurements(ctx context.Context, namespace string) (measurement_repository.MeasurementRepository, error)
	Config(ctx context.Context, namespace string) (config_repository.ConfigRepository, error)
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

func (r *defaultRepository) Config(ctx context.Context, namespace string) (config_repository.ConfigRepository, error) {
	if namespace == "" {
		namespace = "blocks-db"
	}
	collection := r.mongodbClient.Database(namespace).Collection("user_config")
	configRepository, err := config_repository.New(ctx, collection, r.crypto, r.internalEncryptionKeys)
	if err != nil {
		return nil, err
	}
	return configRepository, nil
}

func (r *defaultRepository) Measurements(ctx context.Context, namespace string) (measurement_repository.MeasurementRepository, error) {
	if namespace == "" {
		namespace = "blocks-db"
	}
	db := r.mongodbClient.Database(namespace)
	userActiveMeasurementCollection := db.Collection("user_active_measurements")
	userActiveHistoryCollection := db.Collection("user_active_history")
	namespaceActiveHistoryCollection := db.Collection("namespace_active_history")
	measurementRepository, err := measurement_repository.New(ctx, userActiveMeasurementCollection, userActiveHistoryCollection, namespaceActiveHistoryCollection)
	if err != nil {
		return nil, err
	}
	return measurementRepository, nil
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
