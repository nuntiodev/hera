package config_repository

import (
	"context"
	"errors"
	"github.com/nuntiodev/block-proto/go_block"
	"github.com/nuntiodev/nuntio-user-block/models"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

func (c *defaultConfigRepository) UpdateSettings(ctx context.Context, config *go_block.Config) (*models.Config, error) {
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
			"login_type":                        config.LoginType,
			"require_phone_number_verification": config.RequireEmailVerification,
			"updated_at":                        time.Now(),
		},
	}
	var resp models.Config
	if err := c.collection.FindOneAndUpdate(ctx, bson.M{"_id": namespaceConfigName}, mongoUpdate).Decode(&resp); err != nil {
		return nil, err
	}
	if err := c.crypto.Decrypt(&resp); err != nil {
		return nil, err
	}
	// set updated fields
	resp.EnableNuntioConnect = config.EnableNuntioConnect
	resp.DisableDefaultSignup = config.DisableDefaultSignup
	resp.DisableDefaultLogin = config.DisableDefaultLogin
	resp.ValidatePassword = config.ValidatePassword
	resp.RequireEmailVerification = config.RequireEmailVerification
	resp.LoginType = config.LoginType
	resp.RequirePhoneNumberVerification = config.RequirePhoneNumberVerification
	resp.UpdatedAt = time.Now()
	return &resp, nil
}
