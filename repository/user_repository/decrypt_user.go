package user_repository

import (
	"context"
	"errors"
)

func (r *mongodbRepository) decryptUser(ctx context.Context, user *User, upgrade bool) error {
	if user == nil {
		return errors.New("user is nil")
	}
	// encrypt using internal keys first
	if user.InternalEncrypted {
		if len(r.internalEncryptionKeys) >= user.InternalEncryptionLevel {
			internalKey, err := r.crypto.CombineSymmetricKeys(r.internalEncryptionKeys, user.InternalEncryptionLevel)
			if err != nil {
				return err
			}
			if err := r.decrypt(user, internalKey); err != nil {
				return err
			}
		} else {
			return errors.New("missing required internal encryption keys")
		}
	}
	if user.ExternalEncrypted && r.externalEncryptionKey != "" {
		if err := r.decrypt(user, r.externalEncryptionKey); err != nil {
			return err
		}
	}
	return nil
}

func (r *mongodbRepository) decrypt(user *User, encryptionKey string) error {
	if user == nil {
		return errors.New("user is nil")
	}
	if user.Email != "" {
		decEmail, err := r.crypto.Decrypt(user.Email, encryptionKey)
		if err != nil {
			return err
		}
		user.Email = decEmail
	}
	if user.Image != "" {
		decImage, err := r.crypto.Decrypt(user.Image, encryptionKey)
		if err != nil {
			return err
		}
		user.Image = decImage
	}
	if user.Metadata != "" {
		decMetadata, err := r.crypto.Decrypt(user.Metadata, encryptionKey)
		if err != nil {
			return err
		}
		user.Metadata = decMetadata
	}
	return nil
}
