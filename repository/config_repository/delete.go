package config_repository

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
)

func (c *defaultConfigRepository) Delete(ctx context.Context) error {
	filter := bson.M{"_id": namespaceConfigName}
	result, err := c.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("no documents deleted")
	}
	return nil
}
