package config_repository

import (
	"context"
	"errors"
	"github.com/nuntiodev/block-proto/go_block"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

func (c *defaultConfigRepository) UpdateSettings(ctx context.Context, config *go_block.Config) (*go_block.Config, error) {
	if config == nil {
		return nil, errors.New("missing required config")
	} else if config.Id == "" {
		return nil, errors.New("missing required config")
	}
	get, err := c.GetNamespaceConfig(ctx)
	if err != nil {
		return nil, err
	}
	mongoUpdate := bson.M{
		"$set": bson.M{
			"enable_nuntio_connect":  config.EnableNuntioConnect,
			"disable_default_signup": config.DisableDefaultSignup,
			"disable_default_login":  config.DisableDefaultLogin,
			"validate_password":      config.ValidatePassword,
			"updated_at":             time.Now(),
		},
	}
	if _, err := c.collection.UpdateOne(ctx, bson.M{"_id": namespaceConfigName}, mongoUpdate); err != nil {
		return nil, err
	}
	// set updated fields
	get.EnableNuntioConnect = config.EnableNuntioConnect
	get.DisableDefaultSignup = config.DisableDefaultSignup
	get.DisableDefaultLogin = config.DisableDefaultLogin
	get.ValidatePassword = config.ValidatePassword
	return get, nil
}
