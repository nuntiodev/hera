package crypto

import (
	"errors"
	"fmt"
	"github.com/softcorp-io/block-proto/go_block"
)

func (c *defaultCrypto) DecryptUser(key string, user *go_block.User) error {
	if user == nil {
		return errors.New("user is nil")
	}
	if key == "" {
		return errors.New("key is empty")
	}
	if user.Email != "" {
		decEmail, err := c.Decrypt(user.Email, key)
		if err != nil {
			fmt.Println("2", "email is here", user.Email)
			return err
		}
		user.Email = decEmail
	}
	if user.Image != "" {
		decImage, err := c.Decrypt(user.Image, key)
		if err != nil {
			fmt.Println("3")
			return err
		}
		user.Image = decImage
	}
	if user.Metadata != "" {
		decMetadata, err := c.Decrypt(user.Metadata, key)
		if err != nil {
			fmt.Println("4")
			return err
		}
		user.Metadata = decMetadata
	}
	return nil
}
