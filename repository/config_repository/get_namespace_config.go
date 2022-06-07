package config_repository

import (
	"context"
	"github.com/nuntiodev/nuntio-user-block/models"
	"go.mongodb.org/mongo-driver/bson"
)

func (c *defaultConfigRepository) GetNamespaceConfig(ctx context.Context) (*models.Config, error) {
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
	return &resp, nil
}
