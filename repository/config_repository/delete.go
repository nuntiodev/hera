package config_repository

import (
	"context"
	"errors"
	"github.com/nuntiodev/block-proto/go_block"
	"go.mongodb.org/mongo-driver/bson"
)

func (r *defaultConfigRepository) Delete(ctx context.Context, config *go_block.Config) error {
	if config == nil {
		return errors.New("missing required config")
	} else if config.Id == "" {
		return errors.New("missing required id")
	}
	filter := bson.M{"_id": config.Id}
	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("no documents deleted")
	}
	return nil
}
