package config_repository

import (
	"context"
	"errors"
	"github.com/nuntiodev/block-proto/go_block"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

func (cr *defaultConfigRepository) UpdateLoginText(ctx context.Context, config *go_block.Config) (*go_block.Config, error) {
	if config == nil {
		return nil, errors.New("missing required config")
	} else if config.Id == "" {
		return nil, errors.New("missing required config id")
	}
	get, err := cr.Get(ctx, config)
	if err != nil {
		return nil, err
	}
	update := ProtoConfigToConfig(config)
	if get.InternalEncryptionLevel > 0 {
		update.InternalEncryptionLevel = get.InternalEncryptionLevel
		if err := cr.EncryptConfig(actionUpdate, update); err != nil {
			return nil, err
		}
	}
	updateLoginText := bson.M{}
	if config.LoginText != nil {
		updateLoginText = bson.M{
			"login_button":    update.LoginText.LoginButton,
			"login_title":     update.LoginText.LoginTitle,
			"login_details":   update.LoginText.LoginDetails,
			"forgot_password": update.LoginText.ForgotPassword,
		}
	}
	mongoUpdate := bson.M{
		"$set": bson.M{
			"login_text": updateLoginText,
			"updated_at": time.Now(),
		},
	}
	if _, err := cr.collection.UpdateOne(ctx, bson.M{"_id": config.Id}, mongoUpdate); err != nil {
		return nil, err
	}
	// set updated fields
	get.WelcomeText = config.WelcomeText
	return get, nil
}
