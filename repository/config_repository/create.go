package config_repository

import (
	"context"
	"errors"
	"github.com/nuntiodev/block-proto/go_block"
	ts "google.golang.org/protobuf/types/known/timestamppb"
)

func (c *defaultConfigRepository) Create(ctx context.Context, config *go_block.Config) (*go_block.Config, error) {
	prepare(actionCreate, config)
	if config == nil {
		return nil, errors.New("missing required config")
	}
	// set default fields
	config.EnableNuntioConnect = false
	config.DisableDefaultSignup = false
	config.DisableDefaultLogin = false
	config.ValidatePassword = true
	config.RequireEmailVerification = true
	config.CreatedAt = ts.Now()
	config.UpdatedAt = ts.Now()
	config.DefaultLanguage = go_block.LanguageCode_EN
	config.LoginType = go_block.LoginType_LOGIN_TYPE_EMAIL_PASSWORD
	config.RequirePhoneNumberVerification = true
	config.Id = namespaceConfigName
	create := ProtoConfigToConfig(config)
	if len(c.internalEncryptionKeys) > 0 {
		if err := c.EncryptConfig(actionCreate, create); err != nil {
			return nil, err
		}
		create.InternalEncryptionLevel = int32(len(c.internalEncryptionKeys))
	}
	if _, err := c.collection.InsertOne(ctx, create); err != nil {
		return nil, err
	}
	// set created fields
	config.InternalEncryptionLevel = create.InternalEncryptionLevel
	return config, nil
}
