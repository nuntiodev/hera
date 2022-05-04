package config_repository

import (
	"context"
	"errors"
	"github.com/nuntiodev/block-proto/go_block"
	"go.mongodb.org/mongo-driver/bson"
)

func (c *defaultConfigRepository) GetNamespaceConfig(ctx context.Context) (*go_block.Config, error) {
	resp := Config{}
	result := c.collection.FindOne(ctx, bson.M{"_id": namespaceConfigName})
	if err := result.Err(); err != nil {
		return nil, err
	}
	if err := result.Decode(&resp); err != nil {
		return nil, err
	}
	if resp.InternalEncryptionLevel > 0 && len(c.internalEncryptionKeys) > 0 {
		if resp.InternalEncryptionLevel > int32(len(c.internalEncryptionKeys)) {
			return nil, errors.New("internal encryption level is illegally higher than amount of internal encryption keys")
		}
		if err := c.DecryptConfig(&resp); err != nil {
			return nil, err
		}
		if resp.InternalEncryptionLevel > int32(len(c.internalEncryptionKeys)) {
			// upgrade user to new internal encryption level
			if err := c.upgradeEncryptionLevel(ctx, resp); err != nil {
				return nil, err
			}
		}
	}
	return ConfigToProtoConfig(&resp), nil
}
