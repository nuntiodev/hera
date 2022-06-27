package user_repository

import (
	"context"
	"errors"
	"github.com/nuntiodev/hera-proto/go_hera"
	"github.com/nuntiodev/hera/models"
	"github.com/nuntiodev/x/cryptox"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"time"
)

const (
	actionCreate = iota
	actionGet
	actionUpdateProfile
	actionUpdateMetadata
	actionUpdateContact
	actionUpdatePassword
	actionUpdateEmailVerificationCode
	actionUpdateResetPasswordEmailSent
	actionUpdateTextVerificationCode
)

const (
	maximumGetLimit = 75
	emailHashIndex  = "hera_email_hash_index"
	usernameIndex   = "hera_username_index"
	phoneIndex      = "hera_phone_index"
)

var (
	NoUsersDeletedErr  = errors.New("no users deleted")
	NothingToUpdateErr = errors.New("nothing to update")
	UserIsNilErr       = errors.New("user is nil")
	UpdateIsNil        = errors.New("update is nil")
)

type UserRepository interface {
	Create(ctx context.Context, user *go_hera.User) (*models.User, error)
	UpdateMetadata(ctx context.Context, get *go_hera.User, update *go_hera.User) error
	UpdateProfile(ctx context.Context, get *go_hera.User, update *go_hera.User) error
	UpdatePassword(ctx context.Context, get *go_hera.User, update *go_hera.User) error
	UpdateContact(ctx context.Context, get *go_hera.User, update *go_hera.User) error
	VerifyEmail(ctx context.Context, user *go_hera.User, isVerified bool) error
	UpdateEmailVerificationCode(ctx context.Context, get *go_hera.User) error
	VerifyPhone(ctx context.Context, user *go_hera.User, isVerified bool) error
	UpdatePhoneVerificationCode(ctx context.Context, get *go_hera.User) error
	UpdateResetPasswordCode(ctx context.Context, user *go_hera.User) error
	Get(ctx context.Context, user *go_hera.User) (*models.User, error)
	GetMany(ctx context.Context, users []*go_hera.User) ([]*models.User, error)
	Search(ctx context.Context, search string) (*models.User, error)
	List(ctx context.Context, query *go_hera.Query) ([]*models.User, error)
	Count(ctx context.Context) (int64, error)
	Delete(ctx context.Context, user *go_hera.User) error
	DeleteMany(ctx context.Context, userBatch []*go_hera.User) error
	DeleteAll(ctx context.Context) error
	BuildIndexes(ctx context.Context) error
}

type mongodbRepository struct {
	collection             *mongo.Collection
	crypto                 cryptox.Crypto
	logger                 *zap.Logger
	validatePassword       bool
	maxCodeVerificationAge time.Duration
}

func (r *mongodbRepository) BuildIndexes(ctx context.Context) error {
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
	if _, err := r.collection.Indexes().CreateOne(ctx, emailNamespaceIndexModel); err != nil {
		return err
	}
	usernameIndexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "username_hash", Value: 1},
		},
		Options: options.Index().SetUnique(true).SetPartialFilterExpression(
			bson.D{
				{
					"username_hash", bson.D{
						{
							"$gt", "",
						},
					},
				},
			},
		).SetName(usernameIndex),
	}
	if _, err := r.collection.Indexes().CreateOne(ctx, usernameIndexModel); err != nil {
		return err
	}
	phoneIndexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "phone_hash", Value: 1},
		},
		Options: options.Index().SetUnique(true).SetPartialFilterExpression(
			bson.D{
				{
					"phone_hash", bson.D{
						{
							"$gt", "",
						},
					},
				},
			},
		).SetName(phoneIndex),
	}
	if _, err := r.collection.Indexes().CreateOne(ctx, phoneIndexModel); err != nil {
		return err
	}
	return nil
}

func New(collection *mongo.Collection, crypto cryptox.Crypto, validatePassword bool, maxEmailVerificationAge time.Duration) (UserRepository, error) {
	return &mongodbRepository{
		collection:             collection,
		crypto:                 crypto,
		validatePassword:       validatePassword,
		maxCodeVerificationAge: maxEmailVerificationAge,
	}, nil
}
