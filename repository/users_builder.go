package repository

import (
	"context"
	"github.com/nuntiodev/nuntio-user-block/repository/user_repository"
	"github.com/nuntiodev/x/cryptox"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type UsersBuilder interface {
	SetNamespace(namespace string) UsersBuilder
	SetEncryptionKey(encryptionKey string) UsersBuilder
	WithPasswordValidation(validatePassword bool) UsersBuilder
	Build(ctx context.Context) (user_repository.UserRepository, error)
}

type usersBuilder struct {
	namespace               string
	externalEncryptionKey   string
	validatePassword        bool
	internalEncryptionKeys  []string
	client                  *mongo.Client
	crypto                  cryptox.Crypto
	maxEmailVerificationAge time.Duration
}

func (ub *usersBuilder) SetNamespace(namespace string) UsersBuilder {
	ub.namespace = namespace
	return ub
}

func (ub *usersBuilder) SetEncryptionKey(encryptionKey string) UsersBuilder {
	ub.externalEncryptionKey = encryptionKey
	return ub
}

func (ub *usersBuilder) WithPasswordValidation(validatePassword bool) UsersBuilder {
	ub.validatePassword = validatePassword
	return ub
}

func (ub *usersBuilder) Build(ctx context.Context) (user_repository.UserRepository, error) {
	if ub.namespace == "" {
		ub.namespace = "nuntio-blocks-db"
	}
	collection := ub.client.Database(ub.namespace).Collection("users")
	userRepository, err := user_repository.New(ctx, collection, ub.crypto, ub.internalEncryptionKeys, ub.externalEncryptionKey, ub.validatePassword, ub.maxEmailVerificationAge)
	if err != nil {
		return nil, err
	}
	return userRepository, nil
}

func (r *defaultRepository) Users() UsersBuilder {
	return &usersBuilder{
		crypto:                  r.crypto,
		client:                  r.mongodbClient,
		internalEncryptionKeys:  r.internalEncryptionKeys,
		maxEmailVerificationAge: r.maxEmailVerificationAge,
	}
}
