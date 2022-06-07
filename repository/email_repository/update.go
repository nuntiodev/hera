package email_repository

import (
	"context"
	"errors"
	"github.com/nuntiodev/block-proto/go_block"
	"github.com/nuntiodev/nuntio-user-block/models"
	"github.com/nuntiodev/x/cryptox"
	"go.mongodb.org/mongo-driver/bson"
)

func (e *defaultEmailRepository) Update(ctx context.Context, email *go_block.Email) (*models.Email, error) {
	if email == nil {
		return nil, errors.New("get is nil")
	} else if email.Id == "" {
		return nil, errors.New("missing required id")
	}
	prepare(actionUpdate, email)
	update := models.ProtoEmailToEmail(&go_block.Email{
		Logo:           email.Logo,
		WelcomeMessage: email.WelcomeMessage,
		BodyMessage:    email.BodyMessage,
		FooterMessage:  email.FooterMessage,
		Subject:        email.Subject,
		TemplatePath:   email.TemplatePath,
	})
	if err := e.crypto.Encrypt(update); err != nil {
		return nil, err
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
		},
	}
	result := e.collection.FindOneAndUpdate(ctx, bson.M{"_id": email.Id}, mongoUpdate)
	if err := result.Err(); err != nil {
		return nil, err
	}
	var resp models.Email
	if err := result.Decode(&resp); err != nil {
		return nil, err
	}
	if err := e.crypto.Decrypt(&resp); err != nil {
		return nil, err
	}
	// set updated fields
	resp.Logo = cryptox.Stringx{
		Body:                    email.Logo,
		InternalEncryptionLevel: update.Logo.InternalEncryptionLevel,
		ExternalEncryptionLevel: update.Logo.ExternalEncryptionLevel,
	}
	resp.WelcomeMessage = cryptox.Stringx{
		Body:                    email.WelcomeMessage,
		InternalEncryptionLevel: update.WelcomeMessage.InternalEncryptionLevel,
		ExternalEncryptionLevel: update.WelcomeMessage.ExternalEncryptionLevel,
	}
	resp.BodyMessage = cryptox.Stringx{
		Body:                    email.BodyMessage,
		InternalEncryptionLevel: update.BodyMessage.InternalEncryptionLevel,
		ExternalEncryptionLevel: update.BodyMessage.ExternalEncryptionLevel,
	}
	resp.FooterMessage = cryptox.Stringx{
		Body:                    email.FooterMessage,
		InternalEncryptionLevel: update.FooterMessage.InternalEncryptionLevel,
		ExternalEncryptionLevel: update.FooterMessage.ExternalEncryptionLevel,
	}
	resp.Subject = cryptox.Stringx{
		Body:                    email.Subject,
		InternalEncryptionLevel: update.Subject.InternalEncryptionLevel,
		ExternalEncryptionLevel: update.Subject.ExternalEncryptionLevel,
	}
	resp.TemplatePath = cryptox.Stringx{
		Body:                    email.TemplatePath,
		InternalEncryptionLevel: update.TemplatePath.InternalEncryptionLevel,
		ExternalEncryptionLevel: update.TemplatePath.ExternalEncryptionLevel,
	}
	return &resp, nil
}
