package user_repository

import (
	"errors"
)

type EncryptionOptions struct {
	Key string
}

func (r *mongoRepository) encryptUser(key string, user *User) error {
	if user == nil {
		return errors.New("user is nil")
	}
	if key == "" {
		return errors.New("key is empty")
	}
	if user.Email != "" {
		encEmail, err := r.crypto.Encrypt(user.Email, key)
		if err != nil {
			return err
		}
		user.Email = encEmail
	}
	if user.Image != "" {
		encImage, err := r.crypto.Encrypt(user.Image, key)
		if err != nil {
			return err
		}
		user.Image = encImage
	}
	if user.Metadata != "" {
		encMetadata, err := r.crypto.Encrypt(user.Metadata, key)
		if err != nil {
			return err
		}
		user.Metadata = encMetadata
	}
	return nil
}

func (r *mongoRepository) decryptUser(key string, user *User) error {
	if user == nil {
		return errors.New("user is nil")
	}
	if key == "" {
		return errors.New("key is empty")
	}
	if user.Email != "" {
		decEmail, err := r.crypto.Decrypt(user.Email, key)
		if err != nil {
			return err
		}
		user.Email = decEmail
	}
	if user.Image != "" {
		decImage, err := r.crypto.Decrypt(user.Image, key)
		if err != nil {
			return err
		}
		user.Image = decImage
	}
	if user.Metadata != "" {
		decMetadata, err := r.crypto.Decrypt(user.Metadata, key)
		if err != nil {
			return err
		}
		user.Metadata = decMetadata
	}
	return nil
}
