package email_repository

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/nuntiodev/block-proto/go_block"
	"github.com/nuntiodev/nuntio-user-block/models"
	ts "google.golang.org/protobuf/types/known/timestamppb"
)

func (e *defaultEmailRepository) Create(ctx context.Context, email *go_block.Email) (*models.Email, error) {
	if email == nil {
		return nil, errors.New("email is nil")
	}
	prepare(actionCreate, email)
	if email.Id == "" {
		email.Id = uuid.NewString()
	}
	// set default fields
	email.CreatedAt = ts.Now()
	email.UpdatedAt = ts.Now()
	email.LanguageCode = go_block.LanguageCode_EN
	create := models.ProtoEmailToEmail(email)
	copy := *create
	if err := e.crypto.Encrypt(create); err != nil {
		return nil, err
	}
	if _, err := e.collection.InsertOne(ctx, create); err != nil {
		return nil, err
	}
	// set updated fields
	return &copy, nil
}
