package repository

import (
	"context"
	"github.com/nuntiodev/nuntio-user-block/repository/user_repository"
	"github.com/nuntiodev/x/cryptox"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type UserRepositoryBuilder interface {
	SetNamespace(namespace string) UserRepositoryBuilder
	SetEncryptionKey(encryptionKey string) UserRepositoryBuilder
	WithPasswordValidation(validatePassword bool) UserRepositoryBuilder
	Build(ctx context.Context) (user_repository.UserRepository, error)
}

type usersBuilder struct {
	namespace               string
	externalEncryptionKey   string
	validatePassword        bool
	internalEncryptionKeys  []string
	client                  *mongo.Client
	maxEmailVerificationAge time.Duration
}

func (ub *usersBuilder) SetNamespace(namespace string) UserRepositoryBuilder {
	ub.namespace = namespace
	return ub
}

func (ub *usersBuilder) SetEncryptionKey(encryptionKey string) UserRepositoryBuilder {
	ub.externalEncryptionKey = encryptionKey

	return ub
}

func (ub *usersBuilder) WithPasswordValidation(validatePassword bool) UserRepositoryBuilder {
	ub.validatePassword = validatePassword
	return ub
}

func (ub *usersBuilder) Build(ctx context.Context) (user_repository.UserRepository, error) {
	if ub.namespace == "" {
		ub.namespace = "nuntio-blocks-db"
	}
	crypto, err := cryptox.New(ub.internalEncryptionKeys, []string{ub.externalEncryptionKey})
	if err != nil {
		return nil, err
	}
	collection := ub.client.Database(ub.namespace).Collection("users")
	userRepository, err := user_repository.New(ctx, collection, crypto, ub.validatePassword, ub.maxEmailVerificationAge)
	if err != nil {
		return nil, err
	}
	return userRepository, nil
}

func (r *defaultRepository) UserRepositoryBuilder() UserRepositoryBuilder {
	return &usersBuilder{
		client:                  r.mongodbClient,
		internalEncryptionKeys:  r.internalEncryptionKeys,
		maxEmailVerificationAge: r.maxEmailVerificationAge,
	}
}
