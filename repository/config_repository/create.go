package config_repository

import (
	"context"
	"errors"
	"github.com/nuntiodev/block-proto/go_block"
)

func (cr *defaultConfigRepository) Create(ctx context.Context, config *go_block.Config) (*go_block.Config, error) {
	prepare(actionCreate, config)
	if config == nil {
		return nil, errors.New("missing required config")
	}
	create := ProtoConfigToConfig(config)
	if err := cr.EncryptConfig(actionCreate, create); err != nil {
		return nil, err
	}
	if _, err := cr.collection.InsertOne(ctx, create); err != nil {
		return nil, err
	}
	return config, nil
}
