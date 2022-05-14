package config_repository

import "errors"

func (c *defaultConfigRepository) DecryptConfig(config *Config) error {
	if config == nil {
		return errors.New("decrypt: config is nil")
	}
	if len(c.internalEncryptionKeys) > 0 {
		encryptionKey, err := c.crypto.CombineSymmetricKeys(c.internalEncryptionKeys, len(c.internalEncryptionKeys))
		if err != nil {
			return err
		}
		if config.Name != "" {
			name, err := c.crypto.Decrypt(config.Name, encryptionKey)
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
	}
	return nil
}
