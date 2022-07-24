package repository

import (
	"context"
	"errors"
	"os"
	"strings"
	"time"

	"github.com/nuntiodev/hera-sdks/go_hera"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

const (
	defaultDb = "hera_db"
)

var (
	inMemoryConfigs = []string{}
)

type Repository interface {
	Liveness(ctx context.Context) error
	UserRepositoryBuilder() UserRepositoryBuilder
	TokenRepositoryBuilder() TokenRepositoryBuilder
	ConfigRepositoryBuilder() ConfigRepositoryBuilder
	DropDatabase(ctx context.Context, namespace string) error
	SetDefaultConfig(config *go_hera.Config)
}

type defaultRepository struct {
	namespace               string
	mongodbClient           *mongo.Client
	internalEncryptionKeys  []string
	maxEmailVerificationAge time.Duration
	config                  *go_hera.Config
	logger                  *zap.Logger
	inMemoryConfigs         map[string]*go_hera.Config
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

func (r *defaultRepository) SetDefaultConfig(config *go_hera.Config) {
	r.config = config
}

func initialize() error {
	inMemoryConfigs = strings.Fields(os.Getenv("IN_MEMORY_CONFIGS"))
	for index, config := range inMemoryConfigs {
		inMemoryConfigs[index] = strings.TrimSpace(config)
	}
	return nil
}

func New(mongoClient *mongo.Client, encryptionKeys []string, logger *zap.Logger, maxEmailVerificationAge time.Duration) (Repository, error) {
	logger.Info("creating repository...")
	if err := initialize(); err != nil {
		return nil, err
	}
	repository := &defaultRepository{
		mongodbClient:           mongoClient,
		internalEncryptionKeys:  encryptionKeys,
		maxEmailVerificationAge: maxEmailVerificationAge,
		logger:                  logger,
		inMemoryConfigs:         map[string]*go_hera.Config{},
	}
	for _, config := range inMemoryConfigs {
		configRepository, err := repository.ConfigRepositoryBuilder().SetNamespace(config).Build(context.Background())
		if err != nil {
			return nil, err
		}
		inMemoryConfig, err := configRepository.Get(context.Background())
		if err != nil {
			return nil, err
		}
		repository.inMemoryConfigs[config] = inMemoryConfig
	}
	return repository, nil
}
