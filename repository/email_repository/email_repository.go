package email_repository

import (
	"context"
	"github.com/nuntiodev/block-proto/go_block"
	"github.com/nuntiodev/x/cryptox"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	actionCreate = iota
	actionUpdate
	actionGet
	actionDelete
)

const (
	VerificationEmail  = "verification-email-id"
	ResetPasswordEmail = "reset-password-email-id"
)

type EmailRepository interface {
	Create(ctx context.Context, email *go_block.Email) (*go_block.Email, error)
	Get(ctx context.Context, email *go_block.Email) (*go_block.Email, error)
	GetAll(ctx context.Context, email *go_block.Email) ([]*go_block.Email, error)
	Update(ctx context.Context, email *go_block.Email) (*go_block.Email, error)
	Delete(ctx context.Context, email *go_block.Email) error
}

type defaultEmailRepository struct {
	collection             *mongo.Collection
	crypto                 cryptox.Crypto
	internalEncryptionKeys []string
}

func New(collection *mongo.Collection, crypto cryptox.Crypto, internalEncryptionKeys []string) (EmailRepository, error) {
	return &defaultEmailRepository{
		collection:             collection,
		crypto:                 crypto,
		internalEncryptionKeys: internalEncryptionKeys,
	}, nil
}
