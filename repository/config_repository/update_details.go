package config_repository

import (
	"context"
	"errors"
	"github.com/nuntiodev/block-proto/go_block"
	"github.com/nuntiodev/nuntio-user-block/models"
	"github.com/nuntiodev/x/cryptox"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

func (c *defaultConfigRepository) UpdateDetails(ctx context.Context, config *go_block.Config) (*models.Config, error) {
	if config == nil {
		return nil, errors.New("missing required config")
	} else if config.Id == "" {
		return nil, errors.New("missing required config")
	}
	update := models.ProtoConfigToConfig(&go_block.Config{
		Name:      config.Name,
		Logo:      config.Logo,
		UpdatedAt: config.UpdatedAt,
	})
	if err := c.crypto.Encrypt(update); err != nil {
		return nil, err
	}
	mongoUpdate := bson.M{
		"$set": bson.M{
			"name":       update.Name,
			"logo":       update.Logo,
			"updated_at": time.Now(),
		},
	}
	result := c.collection.FindOneAndUpdate(ctx, bson.M{"_id": namespaceConfigName}, mongoUpdate)
	if err := result.Err(); err != nil {
		return nil, err
	}
	var resp models.Config
	if err := result.Decode(&resp); err != nil {
		return nil, err
	}
	if err := c.crypto.Decrypt(&resp); err != nil {
		return nil, err
	}
	// set updated fields
	resp.Name = cryptox.Stringx{
		Body:                    config.Name,
		InternalEncryptionLevel: update.Name.InternalEncryptionLevel,
		ExternalEncryptionLevel: update.Name.ExternalEncryptionLevel,
	}
	resp.Logo = cryptox.Stringx{
		Body:                    config.Logo,
		InternalEncryptionLevel: update.Logo.InternalEncryptionLevel,
		ExternalEncryptionLevel: update.Logo.ExternalEncryptionLevel,
	}
	return &resp, nil
}
