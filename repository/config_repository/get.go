package config_repository

import (
	"context"
	"errors"
	"github.com/nuntiodev/block-proto/go_block"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

func (cr *defaultConfigRepository) Get(ctx context.Context, config *go_block.Config) (*go_block.Config, error) {
	prepare(actionCreate, config)
	if config == nil {
		return nil, errors.New("missing required config")
	} else if config.Id == "" {
		return nil, errors.New("missing required config id")
	}
	resp := Config{}
	if err := cr.collection.FindOne(ctx, bson.M{"_id": config.Id}).Decode(&resp); err != nil {
		return nil, err
	}
	if resp.InternalEncryptionLevel > 0 && len(cr.internalEncryptionKeys) > 0 {
		if resp.InternalEncryptionLevel > int32(len(cr.internalEncryptionKeys)) {
			return nil, errors.New("internal encryption level is illegally higher than amount of internal encryption keys")
		}
		if err := cr.DecryptConfig(&resp); err != nil {
			return nil, err
		}
		if resp.InternalEncryptionLevel > int32(len(cr.internalEncryptionKeys)) {
			// upgrade user to new internal encryption level
			if err := cr.upgradeEncryptionLevel(ctx, resp); err != nil {
				return nil, err
			}
		}
	}
	return ConfigToProtoConfig(&resp), nil
}

func (cr *defaultConfigRepository) upgradeEncryptionLevel(ctx context.Context, config Config) error {
	if err := cr.EncryptConfig(actionCreate, &config); err != nil {
		return err
	}
	updateAuthConfig := bson.M{}
	if config.AuthConfig != nil {
		updateAuthConfig = bson.M{
			"welcome_title":                 config.AuthConfig.WelcomeTitle,
			"welcome_details":               config.AuthConfig.WelcomeDetails,
			"login_button":                  config.AuthConfig.LoginButton,
			"login_title":                   config.AuthConfig.LoginTitle,
			"login_details":                 config.AuthConfig.LoginDetails,
			"register_button":               config.AuthConfig.RegisterButton,
			"register_title":                config.AuthConfig.RegisterTitle,
			"register_details":              config.AuthConfig.RegisterDetails,
			"missing_password_title":        config.AuthConfig.MissingPasswordTitle,
			"missing_password_details":      config.AuthConfig.MissingPasswordDetails,
			"missing_email_title":           config.AuthConfig.MissingEmailTitle,
			"missing_email_details":         config.AuthConfig.MissingEmailDetails,
			"password_do_not_match_title":   config.AuthConfig.PasswordDoNotMatchTitle,
			"password_do_not_match_details": config.AuthConfig.MissingPasswordDetails,
			"created_by":                    config.AuthConfig.CreatedBy,
		}
	}
	mongoUpdate := bson.M{
		"$set": bson.M{
			"name":                      config.Name,
			"website":                   config.Website,
			"about":                     config.About,
			"email":                     config.Email,
			"logo":                      config.Logo,
			"auth_config":               updateAuthConfig,
			"internal_encryption_level": int32(len(cr.internalEncryptionKeys)),
			"updated_at":                time.Now(),
		},
	}
	if _, err := cr.collection.UpdateOne(ctx, bson.M{"_id": config.Id}, mongoUpdate); err != nil {
		return err
	}
	return nil
}
