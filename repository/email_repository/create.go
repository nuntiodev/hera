package email_repository

import (
	"context"
	"errors"
	"github.com/nuntiodev/block-proto/go_block"
	uuid "github.com/satori/go.uuid"
	ts "google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

func (e *defaultEmailRepository) Create(ctx context.Context, email *go_block.Email) (*go_block.Email, error) {
	if email == nil {
		return nil, errors.New("email is nil")
	}
	prepare(actionCreate, email)
	if email.Id == "" {
		email.Id = uuid.NewV4().String()
	}
	email.CreatedAt = ts.Now()
	email.UpdatedAt = ts.Now()
	create := ProtoEmailToEmail(email)
	if len(e.internalEncryptionKeys) > 0 {
		if err := e.EncryptEmail(actionCreate, create); err != nil {
			return nil, err
		}
		create.InternalEncryptionLevel = int32(len(e.internalEncryptionKeys))
		create.EncryptedAt = time.Now().UTC()
	}
	if _, err := e.collection.InsertOne(ctx, create); err != nil {
		return nil, err
	}
	// set updated fields
	email.InternalEncryptionLevel = create.InternalEncryptionLevel
	email.EncryptedAt = ts.New(create.EncryptedAt)
	return email, nil
}
