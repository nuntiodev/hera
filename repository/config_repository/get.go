package config_repository

import (
	"context"

	"github.com/nuntiodev/hera-sdks/go_hera"
	"github.com/nuntiodev/hera/models"
	"go.mongodb.org/mongo-driver/bson"
)

func (c *defaultConfigRepository) Get(ctx context.Context) (*go_hera.Config, error) {
	if c.config != nil {
		return c.config, nil
	}
	resp := models.Config{}
	result := c.collection.FindOne(ctx, bson.M{"_id": namespaceConfigName})
	if err := result.Err(); err != nil {
		return nil, err
	}
	if err := result.Decode(&resp); err != nil {
		return nil, err
	}
	if err := c.crypto.Decrypt(&resp); err != nil {
		return nil, err
	}
	// check if we should upgrade the encryption level
	if upgradable, _ := c.crypto.Upgradeble(&resp); upgradable {
		if err := c.upgradeEncryptionLevel(ctx, &resp); err != nil {
			return nil, err
		}
	}
	return models.ConfigToProtoConfig(&resp), nil
}
