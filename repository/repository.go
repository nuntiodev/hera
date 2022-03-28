package repository

import (
	"context"
	"github.com/softcorp-io/block-user-service/crypto"
	"github.com/softcorp-io/block-user-service/repository/token_repository"
	"github.com/softcorp-io/block-user-service/repository/user_repository"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"os"
	"time"
)

var (
	defaultNamespaceEncryptionKey = ""
)

type Repository interface {
	Liveness(ctx context.Context) error
	Users(ctx context.Context, namespace, encryptionKey string) (user_repository.UserRepository, error)
	Tokens(ctx context.Context, namespace string) (token_repository.TokenRespository, error)
}

type defaultRepository struct {
	namespace            string
	mongoClient          *mongo.Client
	crypto               crypto.Crypto
	accessTokenExpiresAt time.Duration
}

func initialize() error {
	defaultNamespaceEncryptionKey, _ = os.LookupEnv("DEFAULT_NAMESPACE_ENCRYPTION_KEY")
	return nil
}

func (r *defaultRepository) Liveness(ctx context.Context) error {
	if err := r.mongoClient.Ping(ctx, nil); err != nil {
		return err
	}
	return nil
}

func (r *defaultRepository) Users(ctx context.Context, namespace, encryptionKey string) (user_repository.UserRepository, error) {
	if namespace == "" {
		namespace = "blocks-db"
		if encryptionKey == "" {
			encryptionKey = defaultNamespaceEncryptionKey
		}
	}
	collection := r.mongoClient.Database(namespace).Collection("users")
	userRepository, err := user_repository.New(ctx, collection, r.crypto, encryptionKey)
	if err != nil {
		return nil, err
	}
	return userRepository, nil
}

func (r *defaultRepository) Tokens(ctx context.Context, namespace string) (token_repository.TokenRespository, error) {
	if namespace == "" {
		namespace = "blocks-db"
	}
	collection := r.mongoClient.Database(namespace).Collection("user_tokens")
	tokenRepository, err := token_repository.New(ctx, collection, r.accessTokenExpiresAt)
	if err != nil {
		return nil, err
	}
	return tokenRepository, nil
}

func New(mongoClient *mongo.Client, crypto crypto.Crypto, accessTokenExpiresAt time.Duration, zapLog *zap.Logger) (Repository, error) {
	zapLog.Info("creating repository_mock...")
	if err := initialize(); err != nil {
		return nil, err
	}
	repository := &defaultRepository{
		accessTokenExpiresAt: accessTokenExpiresAt,
		mongoClient:          mongoClient,
		crypto:               crypto,
	}
	return repository, nil
}
