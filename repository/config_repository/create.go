package config_repository

import (
	"context"
	"errors"
	"github.com/nuntiodev/block-proto/go_block"
	"github.com/nuntiodev/nuntio-user-block/models"
	ts "google.golang.org/protobuf/types/known/timestamppb"
)

func (c *defaultConfigRepository) Create(ctx context.Context, config *go_block.Config) (*models.Config, error) {
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
	config.LoginType = go_block.LoginType_LOGIN_TYPE_EMAIL_PASSWORD
	config.RequirePhoneNumberVerification = true
	config.Id = namespaceConfigName
	create := models.ProtoConfigToConfig(config)
	copy := *create
	if err := c.crypto.Encrypt(create); err != nil {
		return nil, err
	}
	if _, err := c.collection.InsertOne(ctx, create); err != nil {
		return nil, err
	}
	// set updated fields
	copy.Name.InternalEncryptionLevel = create.Name.InternalEncryptionLevel
	copy.Name.ExternalEncryptionLevel = create.Name.ExternalEncryptionLevel
	copy.Logo.InternalEncryptionLevel = create.Logo.InternalEncryptionLevel
	copy.Logo.ExternalEncryptionLevel = create.Logo.ExternalEncryptionLevel
	return &copy, nil
}
