package email_repository

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

func (e *defaultEmailRepository) upgradeEncryptionLevel(ctx context.Context, email Email) error {
	if len(e.internalEncryptionKeys) <= 0 {
		return errors.New("length of internal encryption keys is 0")
	}
	if err := e.EncryptEmail(actionCreate, &email); err != nil {
		return err
	}
	mongoUpdate := bson.M{
		"$set": bson.M{
			"logo":                      email.Logo,
			"welcome_message":           email.WelcomeMessage,
			"body_message":              email.BodyMessage,
			"footer_message":            email.FooterMessage,
			"subject":                   email.Subject,
			"template_path":             email.TemplatePath,
			"internal_encryption_level": int32(len(e.internalEncryptionKeys)),
			"updated_at":                time.Now().UTC(),
		},
	}
	if _, err := e.collection.UpdateOne(ctx, bson.M{"_id": email.Id}, mongoUpdate); err != nil {
		return err
	}
	return nil
}
