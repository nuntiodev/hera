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
)

var (
	NoUsersDeletedErr = errors.New("no users deleted")
)

type User struct {
	Id          string    `bson:"_id" json:"id"`
	OptionalId  string    `bson:"optional_id" json:"optional_id"`
	Namespace   string    `bson:"namespace" json:"namespace"`
	Email       string    `bson:"email" json:"email"`
	Role        string    `bson:"role" json:"role"`
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
	Create(ctx context.Context, user *go_block.User, encryptionKey string) (*go_block.User, error)
	UpdatePassword(ctx context.Context, get *go_block.User, update *go_block.User) (*go_block.User, error)
	UpdateEmail(ctx context.Context, get *go_block.User, update *go_block.User, encryptionKey string) (*go_block.User, error)
	UpdateOptionalId(ctx context.Context, get *go_block.User, update *go_block.User) (*go_block.User, error)
	UpdateImage(ctx context.Context, get *go_block.User, update *go_block.User, encryptionKey string) (*go_block.User, error)
	UpdateMetadata(ctx context.Context, get *go_block.User, update *go_block.User, encryptionKey string) (*go_block.User, error)
	UpdateSecurity(ctx context.Context, get *go_block.User, update *go_block.User, encryptionKey string) (*go_block.User, error)
	Get(ctx context.Context, user *go_block.User, encryptionKey string) (*go_block.User, error)
	GetAll(ctx context.Context, userFilter *go_block.UserFilter, namespace string, encryptionKey string) ([]*go_block.User, error)
	GetUsersStream(ctx context.Context, namespace string, userBatch []*go_block.User) (*mongo.ChangeStream, error)
	Count(ctx context.Context, namespace string) (int64, error)
	Delete(ctx context.Context, user *go_block.User) error
	DeleteBatch(ctx context.Context, userBatch []*go_block.User, namespace string) error
	DeleteNamespace(ctx context.Context, namespace string) error
}

type mongoRepository struct {
	collection   *mongo.Collection
	crypto       crypto.Crypto
	metadataType go_block.MetadataType
	zapLog       *zap.Logger
}

func NewUserRepository(ctx context.Context, collection *mongo.Collection, crypto crypto.Crypto, metadataType go_block.MetadataType, zapLog *zap.Logger) (UserRepository, error) {
	zapLog.Info("creating user repository...")
	collection.Indexes().DropAll(context.Background())
	idNamespaceIndexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "_id", Value: 1},
			{Key: "namespace", Value: 1},
		},
		Options: options.Index().SetUnique(true).SetName("id_index_model"),
	}
	if _, err := collection.Indexes().CreateOne(ctx, idNamespaceIndexModel); err != nil {
		return nil, err
	}
	emailNamespaceIndexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "email_hash", Value: 1},
			{Key: "namespace", Value: 1},
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
		).SetName("email_index_model"),
	}
	if _, err := collection.Indexes().CreateOne(ctx, emailNamespaceIndexModel); err != nil {
		return nil, err
	}
	optionalIdIndexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "optional_id", Value: 1},
			{Key: "namespace", Value: 1},
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
		).SetName("optional_id_index_model"),
	}
	if _, err := collection.Indexes().CreateOne(ctx, optionalIdIndexModel); err != nil {
		return nil, err
	}
	return &mongoRepository{
		collection:   collection,
		zapLog:       zapLog,
		crypto:       crypto,
		metadataType: metadataType,
	}, nil
}
