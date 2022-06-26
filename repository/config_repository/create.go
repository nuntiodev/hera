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
	} else if config.Name == "" {
		config.Name = "Hera App"
	}
	// set default fields
	config.ValidatePassword = true
	config.CreatedAt = ts.Now()
	config.UpdatedAt = ts.Now()
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
