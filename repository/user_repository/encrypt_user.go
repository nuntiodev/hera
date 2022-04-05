package user_repository

import (
	"context"
	"errors"
	"time"
)

func (r *mongoRepository) encryptUser(ctx context.Context, action int, user *User) error {
	if user == nil {
		return errors.New("user is nil")
	}
	switch action {
	case actionCreate, actionUpgradeEncryption:
		// encrypt using external key first
		if r.externalEncryptionKey != "" {
			if err := r.encrypt(user, r.externalEncryptionKey); err != nil {
				return err
			}
			user.ExternalEncryptionLevel = 1
			user.ExternalEncrypted = true
			user.EncryptedAt = time.Now()
		}
		// then encrypt using internal key
		if len(r.internalEncryptionKeys) > 0 {
			// check if user has been encrypted before
			internalKey, err := r.crypto.CombineSymmetricKeys(r.internalEncryptionKeys, len(r.internalEncryptionKeys))
			if err != nil {
				return err
			}
			if err := r.encrypt(user, internalKey); err != nil {
				return err
			}
			user.InternalEncryptionLevel = len(r.internalEncryptionKeys)
			user.InternalEncrypted = true
			user.EncryptedAt = time.Now()
		}
	case actionUpdateEmail, actionUpdateImage, actionUpdateSecurity, actionUpdateMetadata:
		if user.ExternalEncrypted && r.externalEncryptionKey != "" {
			if err := r.encrypt(user, r.externalEncryptionKey); err != nil {
				return err
			}
			user.EncryptedAt = time.Now()
		}
		if user.InternalEncrypted && len(r.internalEncryptionKeys) > 0 {
			encryptionKey, err := r.crypto.CombineSymmetricKeys(r.internalEncryptionKeys, user.InternalEncryptionLevel)
			if err != nil {
				return err
			}
			if err := r.encrypt(user, encryptionKey); err != nil {
				return err
			}
			user.EncryptedAt = time.Now()
		}
	default:
		return errors.New("invalid case")
	}
	return nil
}

func (r *mongoRepository) encrypt(user *User, encryptionKey string) error {
	if user == nil {
		return errors.New("user is nil")
	} else if encryptionKey == "" {
		return errors.New("no encryption keys are present")
	}
	if user.Email != "" {
		encEmail, err := r.crypto.Encrypt(user.Email, encryptionKey)
		if err != nil {
			return err
		}
		user.Email = encEmail
	}
	if user.Image != "" {
		encImage, err := r.crypto.Encrypt(user.Image, encryptionKey)
		if err != nil {
			return err
		}
		user.Image = encImage
	}
	if user.Metadata != "" {
		encMetadata, err := r.crypto.Encrypt(user.Metadata, encryptionKey)
		if err != nil {
			return err
		}
		user.Metadata = encMetadata
	}
	user.EncryptedAt = time.Now()
	return nil
}
