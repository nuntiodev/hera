package config_repository

import (
	"context"
	"errors"
	"github.com/nuntiodev/block-proto/go_block"
	"github.com/nuntiodev/nuntio-user-block/models"
	"go.mongodb.org/mongo-driver/bson"
	ts "google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

func (c *defaultConfigRepository) Update(ctx context.Context, config *go_block.Config) (*models.Config, error) {
	if config == nil {
		return nil, errors.New("missing required config")
	} else if config.Id == "" {
		return nil, errors.New("missing required id")
	} else if config.Name == "" {
		return nil, errors.New("missing required name")
	}
	updateConfig := models.ProtoConfigToConfig(&go_block.Config{
		Name:                           config.Name,
		Logo:                           config.Logo,
		EnableNuntioConnect:            config.EnableNuntioConnect,
		DisableDefaultLogin:            config.DisableDefaultLogin,
		DisableDefaultSignup:           config.DisableDefaultSignup,
		ValidatePassword:               config.ValidatePassword,
		RequireEmailVerification:       config.RequireEmailVerification,
		RequirePhoneNumberVerification: config.RequirePhoneNumberVerification,
		UpdatedAt:                      ts.Now(),
	})
	if err := c.crypto.Encrypt(updateConfig); err != nil {
		return nil, err
	}
	mongoUpdate := bson.M{
		"$set": bson.M{
			"name":                       updateConfig.Name,
			"logo":                       updateConfig.Logo,
			"enable_nuntio_connect":      updateConfig.EnableNuntioConnect,
			"disable_default_signup":     updateConfig.DisableDefaultSignup,
			"disable_default_login":      updateConfig.DisableDefaultLogin,
			"validate_password":          updateConfig.ValidatePassword,
			"require_email_verification": updateConfig.RequireEmailVerification,
			//"login_type":                        config.LoginType,
			"require_phone_number_verification": updateConfig.RequireEmailVerification,
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
