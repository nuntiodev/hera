package config_repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

func (c *defaultConfigRepository) RemovePublicKey(ctx context.Context) error {
	mongoUpdate := bson.M{
		"$set": bson.M{
			"public_key": "",
			"updated_at": time.Now(),
		},
	}
	if _, err := c.collection.UpdateOne(ctx, bson.M{"_id": namespaceConfigName}, mongoUpdate); err != nil {
		return err
	}
	return nil
}
