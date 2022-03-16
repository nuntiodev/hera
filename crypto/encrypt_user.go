package crypto

import (
	"errors"
	"github.com/softcorp-io/block-proto/go_block"
)

func (c *defaultCrypto) EncryptUser(key string, user *go_block.User) error {
	if user == nil {
		return errors.New("user is nil")
	}
	if key == "" {
		return errors.New("key is empty")
	}
	if user.Email != "" {
		encEmail, err := c.Encrypt(user.Email, key)
		if err != nil {
			return err
		}
		user.Email = encEmail
	}
	if user.Image != "" {
		encImage, err := c.Encrypt(user.Image, key)
		if err != nil {
			return err
		}
		user.Image = encImage
	}
	if user.Metadata != "" {
		encMetadata, err := c.Encrypt(user.Metadata, key)
		if err != nil {
			return err
		}
		user.Metadata = encMetadata
	}
	return nil
}
