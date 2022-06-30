package repository

import (
	"context"
	"github.com/nuntiodev/hera/repository/user_repository"
	"github.com/nuntiodev/x/cryptox"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type UserRepositoryBuilder interface {
	SetNamespace(namespace string) UserRepositoryBuilder
	WithPasswordValidation(validatePassword bool) UserRepositoryBuilder
	Build(ctx context.Context) (user_repository.UserRepository, error)
}

func NewUserRepositoryBuilder(client *mongo.Client) UserRepositoryBuilder {
	return &userRepositoryBuilder{
		client: client,
	}
}

type userRepositoryBuilder struct {
	namespace               string
	validatePassword        bool
	internalEncryptionKeys  []string
	client                  *mongo.Client
	maxEmailVerificationAge time.Duration
}

func (ub *userRepositoryBuilder) SetNamespace(namespace string) UserRepositoryBuilder {
	ub.namespace = namespace
	return ub
}

func (ub *userRepositoryBuilder) WithPasswordValidation(validatePassword bool) UserRepositoryBuilder {
	ub.validatePassword = validatePassword
	return ub
}

func (ub *userRepositoryBuilder) Build(ctx context.Context) (user_repository.UserRepository, error) {
	if ub.namespace == "" {
		ub.namespace = defaultDb
	}
	crypto, err := cryptox.New(ub.internalEncryptionKeys, nil, nil)
	if err != nil {
		return nil, err
	}
	collection := ub.client.Database(ub.namespace).Collection("users")
	userRepository, err := user_repository.New(collection, crypto, ub.validatePassword, ub.maxEmailVerificationAge)
	if err != nil {
		return nil, err
	}
	return userRepository, nil
}

func (r *defaultRepository) UserRepositoryBuilder() UserRepositoryBuilder {
	return &userRepositoryBuilder{
		client:                  r.mongodbClient,
		internalEncryptionKeys:  r.internalEncryptionKeys,
		maxEmailVerificationAge: r.maxEmailVerificationAge,
	}
}
