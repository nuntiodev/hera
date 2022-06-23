package config_repository

import (
	"context"
	"errors"
	"github.com/nuntiodev/hera-proto/go_hera"
	"github.com/nuntiodev/hera/models"
	ts "google.golang.org/protobuf/types/known/timestamppb"
)

func (c *defaultConfigRepository) Create(ctx context.Context, config *go_hera.Config) error {
	prepare(actionCreate, config)
	if config == nil {
		return errors.New("missing required config")
	}
	// set default fields
	config.DisableSignup = false
	config.DisableLogin = false
	config.ValidatePassword = true
	config.VerifyEmail = true
	config.CreatedAt = ts.Now()
	config.UpdatedAt = ts.Now()
	config.SupportedLoginMechanisms = []go_hera.LoginType{go_hera.LoginType_EMAIL_PASSWORD, go_hera.LoginType_PHONE_PASSWORD, go_hera.LoginType_USERNAME_PASSWORD}
	config.VerifyPhone = true
	create := models.ProtoConfigToConfig(config)
	create.Id = namespaceConfigName
	if err := c.crypto.Encrypt(create); err != nil {
		return err
	}
	if _, err := c.collection.InsertOne(ctx, create); err != nil {
		return err
	}
	return nil
}
