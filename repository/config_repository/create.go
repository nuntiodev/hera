package config_repository

import (
	"context"
	"errors"
	"github.com/nuntiodev/hera-sdks/go_hera"
	"github.com/nuntiodev/hera/models"
	"github.com/nuntiodev/x/randx"
	ts "google.golang.org/protobuf/types/known/timestamppb"
)

func (c *defaultConfigRepository) Create(ctx context.Context, config *go_hera.Config) (*go_hera.Config, error) {
	prepare(actionCreate, config)
	if config == nil {
		return nil, errors.New("missing required config")
	} else if config.Name == "" {
		config.Name = "Hera App"
	}
	// set default fields
	config.ValidatePassword = true
	config.CreatedAt = ts.Now()
	config.UpdatedAt = ts.Now()
	if config.HasingAlgorithm == go_hera.HasingAlgorithm_INVALID_HASHING_ALGORITHM {
		config.HasingAlgorithm = go_hera.HasingAlgorithm_BCRYPT
	}
	if config.Scrypt == nil {
		signerKey, err := randx.GenerateRandomString(88)
		if err != nil {
			return nil, err
		}
		config.Scrypt = &go_hera.Scrypt{SignerKey: signerKey}
	}
	if config.Scrypt.SignerKey == "" {
		return nil, errors.New("missing required signer key for scrypt config")
	}
	create := models.ProtoConfigToConfig(config)
	create.Id = namespaceConfigName
	if err := c.crypto.Encrypt(create); err != nil {
		return nil, err
	}
	if _, err := c.collection.InsertOne(ctx, create); err != nil {
		return nil, err
	}
	return config, nil
}
