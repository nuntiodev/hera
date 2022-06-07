package config_repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/nuntiodev/nuntio-user-block/models"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

func (c *defaultConfigRepository) upgradeEncryptionLevel(ctx context.Context, config *models.Config) error {
	if config == nil {
		return errors.New("config is nil")
	}
	if upgradable, err := c.crypto.Upgradeble(config); err != nil && !upgradable {
		return fmt.Errorf("could not upgrade with err: %v", err)
	}
	copy := *config
	if err := c.crypto.Encrypt(&copy); err != nil {
		return err
	}
	mongoUpdate := bson.M{
		"$set": bson.M{
			"name":       copy.Name,
			"logo":       copy.Logo,
			"updated_at": time.Now(),
		},
	}
	if _, err := c.collection.UpdateOne(ctx, bson.M{"_id": namespaceConfigName}, mongoUpdate); err != nil {
		return err
	}
	config.Name.InternalEncryptionLevel = copy.Name.InternalEncryptionLevel
	config.Name.ExternalEncryptionLevel = copy.Name.ExternalEncryptionLevel
	config.Logo.InternalEncryptionLevel = copy.Logo.InternalEncryptionLevel
	config.Logo.ExternalEncryptionLevel = copy.Logo.ExternalEncryptionLevel
	return nil
}
