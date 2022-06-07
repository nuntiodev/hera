package email_repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/nuntiodev/nuntio-user-block/models"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

func (e *defaultEmailRepository) upgradeEncryptionLevel(ctx context.Context, email *models.Email) error {
	if email == nil {
		return errors.New("email is nil")
	}
	if upgradable, err := e.crypto.Upgradeble(email); err != nil && !upgradable {
		return fmt.Errorf("could not upgrade with err: %v", err)
	}
	copy := *email
	if err := e.crypto.Encrypt(&copy); err != nil {
		return err
	}
	mongoUpdate := bson.M{
		"$set": bson.M{
			"logo":            copy.Logo,
			"welcome_message": copy.WelcomeMessage,
			"body_message":    copy.BodyMessage,
			"footer_message":  copy.FooterMessage,
			"subject":         copy.Subject,
			"template_path":   copy.TemplatePath,
			"updated_at":      time.Now(),
			"encrypted_at":    time.Now(),
		},
	}
	if _, err := e.collection.UpdateOne(
		ctx,
		bson.M{"_id": copy.Id},
		mongoUpdate,
	); err != nil {
		return err
	}
	// update levels
	email.Logo.InternalEncryptionLevel = copy.Logo.InternalEncryptionLevel
	email.Logo.ExternalEncryptionLevel = copy.Logo.ExternalEncryptionLevel
	email.WelcomeMessage.InternalEncryptionLevel = copy.WelcomeMessage.InternalEncryptionLevel
	email.WelcomeMessage.ExternalEncryptionLevel = copy.WelcomeMessage.ExternalEncryptionLevel
	email.BodyMessage.InternalEncryptionLevel = copy.BodyMessage.InternalEncryptionLevel
	email.BodyMessage.ExternalEncryptionLevel = copy.BodyMessage.ExternalEncryptionLevel
	email.FooterMessage.InternalEncryptionLevel = copy.FooterMessage.InternalEncryptionLevel
	email.FooterMessage.ExternalEncryptionLevel = copy.FooterMessage.ExternalEncryptionLevel
	email.Subject.InternalEncryptionLevel = copy.Subject.InternalEncryptionLevel
	email.Subject.ExternalEncryptionLevel = copy.Subject.ExternalEncryptionLevel
	email.TemplatePath.InternalEncryptionLevel = copy.TemplatePath.InternalEncryptionLevel
	email.TemplatePath.ExternalEncryptionLevel = copy.TemplatePath.ExternalEncryptionLevel
	return nil
}
