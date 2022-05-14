package config_repository

import (
	"errors"
)

func (c *defaultConfigRepository) EncryptConfig(action int, config *Config) error {
	if config == nil {
		return errors.New("encrypt: config is nil")
	}
	encryptionKey := ""
	var err error
	switch action {
	case actionCreate:
		encryptionKey, err = c.crypto.CombineSymmetricKeys(c.internalEncryptionKeys, len(c.internalEncryptionKeys))
		if err != nil {
			return err
		}
	case actionUpdate:
		encryptionKey, err = c.crypto.CombineSymmetricKeys(c.internalEncryptionKeys, int(config.InternalEncryptionLevel))
		if err != nil {
			return err
		}
	default:
		return errors.New("invalid action")
	}
	if config.Name != "" {
		name, err := c.crypto.Encrypt(config.Name, encryptionKey)
		if err != nil {
			return err
		}
		config.Name = name
	}
	if config.Logo != "" {
		logo, err := c.crypto.Encrypt(config.Logo, encryptionKey)
		if err != nil {
			return err
		}
		config.Logo = logo
	}
	return nil
}
