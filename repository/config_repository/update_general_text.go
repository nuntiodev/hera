package config_repository

import (
	"context"
	"errors"
	"github.com/nuntiodev/block-proto/go_block"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

func (c *defaultConfigRepository) UpdateGeneralText(ctx context.Context, config *go_block.Config) (*go_block.Config, error) {
	if config == nil {
		return nil, errors.New("missing required config")
	} else if config.Id == "" {
		return nil, errors.New("missing required config id")
	}
	get, err := c.GetNamespaceConfig(ctx)
	if err != nil {
		return nil, err
	}
	update := ProtoConfigToConfig(config)
	if get.InternalEncryptionLevel > 0 {
		update.InternalEncryptionLevel = get.InternalEncryptionLevel
		if err := c.EncryptConfig(actionUpdate, update); err != nil {
			return nil, err
		}
	}
	updateGeneralText := bson.M{}
	if config.GeneralText != nil {
		updateGeneralText = bson.M{
			"missing_password_title":   update.GeneralText.MissingPasswordTitle,
			"missing_password_details": update.GeneralText.MissingPasswordDetails,
			"missing_email_title":      update.GeneralText.MissingEmailTitle,
			"missing_email_details":    update.GeneralText.MissingEmailDetails,
			"created_by":               update.GeneralText.CreatedBy,
			"password_hint":            update.GeneralText.PasswordHint,
			"email_hint":               update.GeneralText.EmailHint,
			"error_title":              update.GeneralText.ErrorTitle,
			"error_description":        update.GeneralText.ErrorDescription,
			"no_wifi_title":            update.GeneralText.NoWifiTitle,
			"no_wifi_description":      update.GeneralText.NoWifiDescription,
		}
	}
	mongoUpdate := bson.M{
		"$set": bson.M{
			"general_text": updateGeneralText,
			"updated_at":   time.Now().UTC(),
		},
	}
	if _, err := c.collection.UpdateOne(ctx, bson.M{"_id": namespaceConfigName}, mongoUpdate); err != nil {
		return nil, err
	}
	// set updated fields
	get.WelcomeText = config.WelcomeText
	return get, nil
}
