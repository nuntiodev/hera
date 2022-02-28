package repository

import (
	"context"
	"errors"
	"github.com/softcorp-io/block-user-service/repository/user_repository"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"os"
)

var (
	mongoName           = ""
	mongoUserCollection = ""
)

type Repository struct {
	UserRepository user_repository.UserRepository
	mongoClient    *mongo.Client
}

func initialize() error {
	var ok bool
	mongoName, ok = os.LookupEnv("MONGO_DB_NAME")
	if !ok || mongoName == "" {
		return errors.New("missing required MONGO_DB_NAME")
	}
	mongoUserCollection, ok = os.LookupEnv("MONGO_USER_COLLECTION")
	if !ok || mongoUserCollection == "" {
		return errors.New("missing required MONGO_USER_COLLECTION")
	}
	return nil
}

func (r *Repository) Liveness(ctx context.Context) error {
	if err := r.mongoClient.Ping(ctx, nil); err != nil {
		return err
	}
	return nil
}

func New(ctx context.Context, mongoClient *mongo.Client, zapLog *zap.Logger) (*Repository, error) {
	zapLog.Info("creating repository_mock...")
	if err := initialize(); err != nil {
		return nil, err
	}
	userCollection := mongoClient.Database(mongoName).Collection(mongoUserCollection)
	userRepository, err := user_repository.NewUserRepository(ctx, userCollection, zapLog)
	if err != nil {
		return nil, err
	}
	repository := &Repository{
		UserRepository: userRepository,
		mongoClient:    mongoClient,
	}
	return repository, nil
}
