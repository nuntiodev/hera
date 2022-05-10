package email_repository

import (
	"context"
	"errors"
	"github.com/nuntiodev/block-proto/go_block"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

func (e *defaultEmailRepository) Update(ctx context.Context, email *go_block.Email) (*go_block.Email, error) {
	if email == nil {
		return nil, errors.New("get is nil")
	} else if email.Id == "" {
		return nil, errors.New("missing required id")
	}
	prepare(actionUpdate, email)
	get, err := e.Get(ctx, email)
	if err != nil {
		return nil, err
	}
	update := ProtoEmailToEmail(email)
	if get.InternalEncryptionLevel > 0 {
		if err := e.EncryptEmail(actionUpdate, update); err != nil {
			return nil, err
		}
		update.EncryptedAt = time.Now()
	}
	mongoUpdate := bson.M{
		"$set": bson.M{
			"logo":            update.Logo,
			"welcome_message": update.WelcomeMessage,
			"body_message":    update.BodyMessage,
			"footer_message":  update.FooterMessage,
			"subject":         update.Subject,
			"template_path":   update.TemplatePath,
			"updated_at":      update.UpdatedAt,
			"encrypted_at":    update.EncryptedAt,
		},
	}
	if _, err := e.collection.UpdateOne(ctx, bson.M{"_id": get.Id}, mongoUpdate); err != nil {
		return nil, err
	}
	// set updated fields
	get.Logo = email.Logo
	get.WelcomeMessage = email.WelcomeMessage
	get.BodyMessage = email.BodyMessage
	get.FooterMessage = email.FooterMessage
	get.Subject = email.Subject
	get.UpdatedAt = email.UpdatedAt
	get.EncryptedAt = email.EncryptedAt
	return email, nil
}
