package config_repository

import (
	"context"
	"errors"
	"github.com/nuntiodev/block-proto/go_block"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

func (cr *defaultConfigRepository) UpdateAuthConfig(ctx context.Context, config *go_block.Config) (*go_block.Config, error) {
	if config == nil {
		return nil, errors.New("missing required config")
	} else if config.Id == "" {
		return nil, errors.New("missing required config id")
	}
	if _, err := cr.Create(ctx, &go_block.Config{}); err != nil {
		return nil, err
	}
	get, err := cr.Get(ctx, config)
	if err != nil {
		return nil, err
	}
	update := ProtoConfigToConfig(config)
	if get.InternalEncryptionLevel > 0 {
		if err := cr.EncryptConfig(actionUpdate, update); err != nil {
			return nil, err
		}
	}
	updateAuthConfig := bson.M{}
	if config.AuthConfig != nil {
		updateAuthConfig = bson.M{
			"welcome_title":                 update.AuthConfig.WelcomeTitle,
			"welcome_details":               update.AuthConfig.WelcomeDetails,
			"login_button":                  update.AuthConfig.LoginButton,
			"login_title":                   update.AuthConfig.LoginTitle,
			"login_details":                 update.AuthConfig.LoginDetails,
			"register_button":               update.AuthConfig.RegisterButton,
			"register_title":                update.AuthConfig.RegisterTitle,
			"register_details":              update.AuthConfig.RegisterDetails,
			"missing_password_title":        update.AuthConfig.MissingPasswordTitle,
			"missing_password_details":      update.AuthConfig.MissingPasswordDetails,
			"missing_email_title":           update.AuthConfig.MissingEmailTitle,
			"missing_email_details":         update.AuthConfig.MissingEmailDetails,
			"password_do_not_match_title":   update.AuthConfig.PasswordDoNotMatchTitle,
			"password_do_not_match_details": update.AuthConfig.MissingPasswordDetails,
			"created_by":                    update.AuthConfig.CreatedBy,
		}
	}
	mongoUpdate := bson.M{
		"$set": bson.M{
			"auth_config": updateAuthConfig,
			"updated_at":  time.Now(),
		},
	}
	if _, err := cr.collection.UpdateOne(ctx, bson.M{"_id": config.Id}, mongoUpdate); err != nil {
		return nil, err
	}
	// set updated fields
	get.AuthConfig = config.AuthConfig
	return get, nil
}
