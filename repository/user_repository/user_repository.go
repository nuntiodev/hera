package user_repository

import (
	"context"
	"errors"
	"github.com/softcorp-io/block-proto/go_block"
	"github.com/softcorp-io/block-user-service/crypto"
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
	actionUpdateOptionalId
	actionUpdateImage
	actionUpdateMetadata
	actionUpdateNamespace
	actionUpdateSecurity
	actionGet
	actionGetAll
)

const (
	maximumGetLimit = 75
	maxFieldLength  = 150
	emailHashIndex  = "block_email_hash_index"
	optionalIdIndex = "block_optional_id_index"
)

var (
	NoUsersDeletedErr = errors.New("no users deleted")
)

type User struct {
	Id          string    `bson:"_id" json:"id"`
	OptionalId  string    `bson:"optional_id" json:"optional_id"`
	Email       string    `bson:"email" json:"email"`
	Password    string    `bson:"password" json:"password"`
	Image       string    `bson:"image" json:"image"`
	Encrypted   bool      `bson:"encrypted" json:"encrypted"`
	EmailHash   string    `bson:"email_hash" json:"email_hash"`
	Metadata    string    `bson:"metadata" json:"metadata"`
	CreatedAt   time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at" json:"updated_at"`
	EncryptedAt time.Time `bson:"encrypted_at" json:"encrypted_at"`
}

type UserRepository interface {
	Create(ctx context.Context, user *go_block.User) (*go_block.User, error)
	UpdatePassword(ctx context.Context, get *go_block.User, update *go_block.User) (*go_block.User, error)
	UpdateEmail(ctx context.Context, get *go_block.User, update *go_block.User) (*go_block.User, error)
	UpdateOptionalId(ctx context.Context, get *go_block.User, update *go_block.User) (*go_block.User, error)
	UpdateImage(ctx context.Context, get *go_block.User, update *go_block.User) (*go_block.User, error)
	UpdateMetadata(ctx context.Context, get *go_block.User, update *go_block.User) (*go_block.User, error)
	UpdateSecurity(ctx context.Context, get *go_block.User) (*go_block.User, error)
	Get(ctx context.Context, user *go_block.User) (*go_block.User, error)
	GetAll(ctx context.Context, userFilter *go_block.UserFilter) ([]*go_block.User, error)
	Count(ctx context.Context) (int64, error)
	Delete(ctx context.Context, user *go_block.User) error
	DeleteBatch(ctx context.Context, userBatch []*go_block.User) error
	DeleteAll(ctx context.Context) error
}

type mongoRepository struct {
	collection    *mongo.Collection
	crypto        crypto.Crypto
	zapLog        *zap.Logger
	encryptionKey string
}

func New(ctx context.Context, collection *mongo.Collection, crypto crypto.Crypto, encryptionKey string) (UserRepository, error) {
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
	optionalIdIndexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "optional_id", Value: 1},
		},
		Options: options.Index().SetUnique(true).SetPartialFilterExpression(
			bson.D{
				{
					"optional_id", bson.D{
						{
							"$gt", "",
						},
					},
				},
			},
		).SetName(optionalIdIndex),
	}
	if _, err := collection.Indexes().CreateOne(ctx, optionalIdIndexModel); err != nil {
		return nil, err
	}
	return &mongoRepository{
		collection:    collection,
		crypto:        crypto,
		encryptionKey: encryptionKey,
	}, nil
}
