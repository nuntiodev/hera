package text_repository

import (
	"context"
	"github.com/nuntiodev/block-proto/go_block"
	"github.com/nuntiodev/x/cryptox"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	actionCreate = iota
	actionUpdate
)

type TextRepository interface {
	Create(ctx context.Context, text *go_block.Text) (*go_block.Text, error)
	Get(ctx context.Context, id go_block.LanguageCode) (*go_block.Text, error)
	GetAll(ctx context.Context) ([]*go_block.Text, error)
	UpdateGeneral(ctx context.Context, text *go_block.Text) (*go_block.Text, error)
	UpdateWelcome(ctx context.Context, text *go_block.Text) (*go_block.Text, error)
	UpdateRegister(ctx context.Context, text *go_block.Text) (*go_block.Text, error)
	UpdateLogin(ctx context.Context, text *go_block.Text) (*go_block.Text, error)
	UpdateProfile(ctx context.Context, text *go_block.Text) (*go_block.Text, error)
	Delete(ctx context.Context, id go_block.LanguageCode) error
}

type defaultTextRepository struct {
	collection             *mongo.Collection
	crypto                 cryptox.Crypto
	internalEncryptionKeys []string
}

func newMongodbTextRepository(ctx context.Context, collection *mongo.Collection, crypto cryptox.Crypto, internalEncryptionKeys []string) (*defaultTextRepository, error) {
	return &defaultTextRepository{
		collection:             collection,
		crypto:                 crypto,
		internalEncryptionKeys: internalEncryptionKeys,
	}, nil
}

func New(ctx context.Context, collection *mongo.Collection, crypto cryptox.Crypto, internalEncryptionKeys []string) (TextRepository, error) {
	return newMongodbTextRepository(ctx, collection, crypto, internalEncryptionKeys)
}
