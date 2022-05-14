package config_repository

import (
	"context"
	"errors"
	"github.com/nuntiodev/block-proto/go_block"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

func (c *defaultConfigRepository) UpdateSettings(ctx context.Context, config *go_block.Config) (*go_block.Config, error) {
	if config == nil {
		return nil, errors.New("missing required config")
	} else if config.Id == "" {
		return nil, errors.New("missing required config")
	}
	mongoUpdate := bson.M{
		"$set": bson.M{
			"enable_nuntio_connect":             config.EnableNuntioConnect,
			"disable_default_signup":            config.DisableDefaultSignup,
			"disable_default_login":             config.DisableDefaultLogin,
			"validate_password":                 config.ValidatePassword,
			"require_email_verification":        config.RequireEmailVerification,
			"default_language":                  config.DefaultLanguage,
			"login_type":                        config.LoginType,
			"require_phone_number_verification": config.RequireEmailVerification,
			"updated_at":                        time.Now(),
		},
	}
	var resp Config
	if err := c.collection.FindOneAndUpdate(ctx, bson.M{"_id": namespaceConfigName}, mongoUpdate).Decode(&resp); err != nil {
		return nil, err
	}
	if resp.InternalEncryptionLevel > 0 {
		if err := c.DecryptConfig(&resp); err != nil {
			return nil, err
		}
	}
	get := ConfigToProtoConfig(&resp)
	// set updated fields
	get.EnableNuntioConnect = config.EnableNuntioConnect
	get.DisableDefaultSignup = config.DisableDefaultSignup
	get.DisableDefaultLogin = config.DisableDefaultLogin
	get.ValidatePassword = config.ValidatePassword
	return get, nil
}
