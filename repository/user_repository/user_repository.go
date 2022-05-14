package user_repository

import (
	"context"
	"errors"
	"github.com/nuntiodev/block-proto/go_block"
	"github.com/nuntiodev/x/cryptox"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"time"
)

const (
	actionCreate = iota
	actionUpdatePassword
	actionUpdateEmail
	actionUpdateUsername
	actionUpdateImage
	actionUpdateMetadata
	actionUpdateNamespace
	actionUpdateBirthdate
	actionUpdateName
	actionUpdateSecurity
	actionUpdateEmailVerified
	actionUpdateVerificationEmailSent
	actionUpdateResetPasswordEmailSent
	actionUpdatePreferredLanguage
	actionGet
	actionGetAll
	actionUpgradeEncryption
	actionUpdatePhoneNumber
)

const (
	maximumGetLimit = 75
	maxFieldLength  = 150
	minEntropy      = 60
	emailHashIndex  = "block_email_hash_index"
	usernameIndex   = "block_username_index"
)

var (
	NoUsersDeletedErr = errors.New("no users deleted")
	UserIsNilErr      = errors.New("user is nil")
)

type UserRepository interface {
	Create(ctx context.Context, user *go_block.User) (*go_block.User, error)
	UpdatePassword(ctx context.Context, get *go_block.User, update *go_block.User) (*go_block.User, error)
	UpdateEmail(ctx context.Context, get *go_block.User, update *go_block.User) (*go_block.User, error)
	UpdateUsername(ctx context.Context, get *go_block.User, update *go_block.User) (*go_block.User, error)
	UpdateImage(ctx context.Context, get *go_block.User, update *go_block.User) (*go_block.User, error)
	UpdateMetadata(ctx context.Context, get *go_block.User, update *go_block.User) (*go_block.User, error)
	UpdateName(ctx context.Context, get *go_block.User, update *go_block.User) (*go_block.User, error)
	UpdateBirthdate(ctx context.Context, get *go_block.User, update *go_block.User) (*go_block.User, error)
	UpdatePreferredLanguage(ctx context.Context, get *go_block.User, update *go_block.User) (*go_block.User, error)
	UpdateSecurity(ctx context.Context, get *go_block.User) (*go_block.User, error)
	UpdateVerificationEmailSent(ctx context.Context, get *go_block.User) (*go_block.User, error)
	UpdateResetPasswordEmailSent(ctx context.Context, user *go_block.User) (*go_block.User, error)
	UpdateEmailVerified(ctx context.Context, get *go_block.User, update *go_block.User) (*go_block.User, error)
	UpdatePhoneNumber(ctx context.Context, get *go_block.User, update *go_block.User) (*go_block.User, error)
	Get(ctx context.Context, user *go_block.User, upgrade bool) (*go_block.User, error)
	GetAll(ctx context.Context, userFilter *go_block.UserFilter) ([]*go_block.User, error)
	Count(ctx context.Context) (int64, error)
	Delete(ctx context.Context, user *go_block.User) error
	DeleteBatch(ctx context.Context, userBatch []*go_block.User) error
	DeleteAll(ctx context.Context) error
}

type mongodbRepository struct {
	collection              *mongo.Collection
	crypto                  cryptox.Crypto
	zapLog                  *zap.Logger
	internalEncryptionKeys  []string
	externalEncryptionKey   string
	validatePassword        bool
	maxEmailVerificationAge time.Duration
}

func newMongodbUserRepository(ctx context.Context, collection *mongo.Collection, crypto cryptox.Crypto, internalEncryptionKeys []string, externalEncryptionKey string, validatePassword bool, maxEmailVerificationAge time.Duration) (*mongodbRepository, error) {
	emailNamespaceIndexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "email_hash", Value: 1},
		},
		Options: options.Index().SetUnique(true).SetPartialFilterExpression(
			bson.D{
				{
					"email_hash", bson.D{
						{
							"$gt", "",
						},
					},
				},
			},
		).SetName(emailHashIndex),
	}
	if _, err := collection.Indexes().CreateOne(ctx, emailNamespaceIndexModel); err != nil {
		return nil, err
	}
	usernameIndexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "username", Value: 1},
		},
		Options: options.Index().SetUnique(true).SetPartialFilterExpression(
			bson.D{
				{
					"username", bson.D{
						{
							"$gt", "",
						},
					},
				},
			},
		).SetName(usernameIndex),
	}
	if _, err := collection.Indexes().CreateOne(ctx, usernameIndexModel); err != nil {
		return nil, err
	}
	return &mongodbRepository{
		collection:              collection,
		crypto:                  crypto,
		internalEncryptionKeys:  internalEncryptionKeys,
		externalEncryptionKey:   externalEncryptionKey,
		validatePassword:        validatePassword,
		maxEmailVerificationAge: maxEmailVerificationAge,
	}, nil
}

func New(ctx context.Context, collection *mongo.Collection, crypto cryptox.Crypto, internalEncryptionKeys []string, externalEncryptionKey string, validatePassword bool, maxEmailVerificationAge time.Duration) (UserRepository, error) {
	return newMongodbUserRepository(ctx, collection, crypto, internalEncryptionKeys, externalEncryptionKey, validatePassword, maxEmailVerificationAge)
}
