package config_repository

import (
	"context"
	"errors"
	"github.com/nuntiodev/block-proto/go_block"
	"go.mongodb.org/mongo-driver/bson"
)

func (cr *defaultConfigRepository) Get(ctx context.Context, config *go_block.Config) (*go_block.Config, error) {
	prepare(actionCreate, config)
	if config == nil {
		return nil, errors.New("missing required config")
	} else if config.Id == "" {
		return nil, errors.New("missing required config id")
	}
	resp := Config{}
	result := cr.collection.FindOne(ctx, bson.M{"_id": config.Id})
	if err := result.Err(); err != nil {
		// create one since it does not exist
		// todo: delete - we do not need to create a new one if it does not exist
		create, err := cr.Create(ctx, config)
		if err != nil {
			return nil, err
		}
		return create, nil
	}
	if err := result.Decode(&resp); err != nil {
		return nil, err
	}
	if resp.InternalEncryptionLevel > 0 && len(cr.internalEncryptionKeys) > 0 {
		if resp.InternalEncryptionLevel > int32(len(cr.internalEncryptionKeys)) {
			return nil, errors.New("internal encryption level is illegally higher than amount of internal encryption keys")
		}
		if err := cr.DecryptConfig(&resp); err != nil {
			return nil, err
		}
		if resp.InternalEncryptionLevel > int32(len(cr.internalEncryptionKeys)) {
			// upgrade user to new internal encryption level
			if err := cr.upgradeEncryptionLevel(ctx, resp); err != nil {
				return nil, err
			}
		}
	}
	return ConfigToProtoConfig(&resp), nil
}
