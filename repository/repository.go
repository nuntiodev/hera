package repository

import (
	"context"
	"errors"
	"github.com/nuntiodev/nuntio-user-block/repository/config_repository"
	"github.com/nuntiodev/nuntio-user-block/repository/email_repository"
	"github.com/nuntiodev/nuntio-user-block/repository/measurement_repository"
	"github.com/nuntiodev/nuntio-user-block/repository/token_repository"
	"github.com/nuntiodev/x/cryptox"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"time"
)

const (
	defaultDb = "nuntio-blocks-db"
)

type Repository interface {
	Liveness(ctx context.Context) error
	UserRepositoryBuilder() UserRepositoryBuilder
	Tokens(ctx context.Context, namespace, externalEncryptionKey string) (token_repository.TokenRepository, error)
	Measurements(ctx context.Context, namespace string) (measurement_repository.MeasurementRepository, error)
	Config(ctx context.Context, namespace, externalEncryptionKey string) (config_repository.ConfigRepository, error)
	Email(ctx context.Context, namespace, externalEncryptionKey string) (email_repository.EmailRepository, error)
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

func (r *defaultRepository) Tokens(ctx context.Context, namespace, externalEncryptionKey string) (token_repository.TokenRepository, error) {
	if namespace == "" {
		namespace = defaultDb
	}
	crypto, err := cryptox.New(r.internalEncryptionKeys, []string{externalEncryptionKey})
	if err != nil {
		return nil, err
	}
	collection := r.mongodbClient.Database(namespace).Collection("user_tokens")
	tokenRepository, err := token_repository.New(ctx, collection, crypto, r.internalEncryptionKeys)
	if err != nil {
		return nil, err
	}
	return tokenRepository, nil
}

func (r *defaultRepository) Measurements(ctx context.Context, namespace string) (measurement_repository.MeasurementRepository, error) {
	if namespace == "" {
		namespace = defaultDb
	}
	db := r.mongodbClient.Database(namespace)
	userActiveMeasurementCollection := db.Collection("user_active_measurements")
	userActiveHistoryCollection := db.Collection("user_active_history")
	measurementRepository, err := measurement_repository.New(ctx, userActiveMeasurementCollection, userActiveHistoryCollection)
	if err != nil {
		return nil, err
	}
	return measurementRepository, nil
}

func (r *defaultRepository) Config(ctx context.Context, namespace, externalEncryptionKey string) (config_repository.ConfigRepository, error) {
	if namespace == "" {
		namespace = defaultDb
	}
	crypto, err := cryptox.New(r.internalEncryptionKeys, []string{externalEncryptionKey})
	if err != nil {
		return nil, err
	}
	collection := r.mongodbClient.Database(namespace).Collection("user_config")
	configRepository, err := config_repository.New(ctx, collection, crypto)
	if err != nil {
		return nil, err
	}
	return configRepository, nil
}

func (r *defaultRepository) Email(ctx context.Context, namespace, externalEncryptionKey string) (email_repository.EmailRepository, error) {
	if namespace == "" {
		namespace = defaultDb
	}
	crypto, err := cryptox.New(r.internalEncryptionKeys, []string{externalEncryptionKey})
	if err != nil {
		return nil, err
	}
	collection := r.mongodbClient.Database(namespace).Collection("user_emails")
	emailRepository, err := email_repository.New(collection, crypto)
	if err != nil {
		return nil, err
	}
	return emailRepository, nil
}

func (r *defaultRepository) DropDatabase(ctx context.Context, namespace string) error {
	if namespace == "" {
		return errors.New("missing required namespace")
	}
	return r.mongodbClient.Database(namespace).Drop(ctx)
}

func New(mongoClient *mongo.Client, encryptionKeys []string, zapLog *zap.Logger, maxEmailVerificationAge time.Duration) (Repository, error) {
	zapLog.Info("creating repository...")
	repository := &defaultRepository{
		mongodbClient:           mongoClient,
		internalEncryptionKeys:  encryptionKeys,
		maxEmailVerificationAge: maxEmailVerificationAge,
	}
	return repository, nil
}
