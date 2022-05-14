package config_repository

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

func (c *defaultConfigRepository) upgradeEncryptionLevel(ctx context.Context, config Config) error {
	if len(c.internalEncryptionKeys) <= 0 {
		return errors.New("length of internal encryption keys is 0")
	}
	if err := c.EncryptConfig(actionCreate, &config); err != nil {
		return err
	}
	mongoUpdate := bson.M{
		"$set": bson.M{
			"name":                      config.Name,
			"logo":                      config.Logo,
			"internal_encryption_level": int32(len(c.internalEncryptionKeys)),
			"updated_at":                time.Now(),
		},
	}
	if _, err := c.collection.UpdateOne(ctx, bson.M{"_id": namespaceConfigName}, mongoUpdate); err != nil {
		return err
	}
	return nil
}
