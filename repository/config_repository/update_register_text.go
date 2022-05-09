package config_repository

import (
	"context"
	"errors"
	"github.com/nuntiodev/block-proto/go_block"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

func (c *defaultConfigRepository) UpdateRegisterText(ctx context.Context, config *go_block.Config) (*go_block.Config, error) {
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
	updateRegisterText := bson.M{}
	if config.RegisterText != nil {
		updateRegisterText = bson.M{
			"register_button":               update.RegisterText.RegisterButton,
			"register_title":                update.RegisterText.RegisterTitle,
			"register_details":              update.RegisterText.RegisterDetails,
			"password_do_not_match_title":   update.RegisterText.PasswordDoNotMatchTitle,
			"password_do_not_match_details": update.RegisterText.PasswordDoNotMatchDetails,
			"repeat_password_hint":          update.RegisterText.RepeatPasswordHint,
			"contains_special_char":         update.RegisterText.ContainsSpecialChar,
			"contains_number_char":          update.RegisterText.ContainsNumberChar,
			"password_must_match":           update.RegisterText.PasswordMustMatch,
			"contains_eight_chars":          update.RegisterText.ContainsEightChars,
		}
	}
	mongoUpdate := bson.M{
		"$set": bson.M{
			"register_text": updateRegisterText,
			"updated_at":    time.Now().UTC(),
		},
	}
	if _, err := c.collection.UpdateOne(ctx, bson.M{"_id": namespaceConfigName}, mongoUpdate); err != nil {
		return nil, err
	}
	// set updated fields
	get.WelcomeText = config.WelcomeText
	return get, nil
}
