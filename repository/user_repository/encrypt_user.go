package user_repository

import (
	"context"
	"errors"
	"fmt"
	"time"
)

func (r *mongodbRepository) encryptUser(ctx context.Context, action int, user *User) error {
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
			user.EncryptedAt = time.Now().UTC()
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
			user.EncryptedAt = time.Now().UTC()
		}
	case actionUpdateEmail, actionUpdateImage, actionUpdateSecurity, actionUpdateMetadata, actionUpdateName, actionUpdateBirthdate:
		if user.ExternalEncryptionLevel > 0 && r.externalEncryptionKey != "" {
			if err := r.encrypt(user, r.externalEncryptionKey); err != nil {
				return err
			}
			user.EncryptedAt = time.Now().UTC()
		}
		if user.InternalEncryptionLevel > 0 && len(r.internalEncryptionKeys) > 0 {
			encryptionKey, err := r.crypto.CombineSymmetricKeys(r.internalEncryptionKeys, user.InternalEncryptionLevel)
			if err != nil {
				return err
			}
			if err := r.encrypt(user, encryptionKey); err != nil {
				return err
			}
			user.EncryptedAt = time.Now().UTC()
		}
	default:
		return errors.New("invalid case")
	}
	return nil
}

func (r *mongodbRepository) encrypt(user *User, encryptionKey string) error {
	if user == nil {
		return errors.New("user is nil")
	} else if encryptionKey == "" {
		return errors.New("no encryption keys are present")
	}
	if user.Email != "" {
		encEmail, err := r.crypto.Encrypt(user.Email, encryptionKey)
		if err != nil {
			return fmt.Errorf("cannot encrypt email: %v", err)
		}
		user.Email = encEmail
	}
	if user.Image != "" {
		encImage, err := r.crypto.Encrypt(user.Image, encryptionKey)
		if err != nil {
			return fmt.Errorf("cannot encrypt image: %v", err)
		}
		user.Image = encImage
	}
	if user.Metadata != "" {
		encMetadata, err := r.crypto.Encrypt(user.Metadata, encryptionKey)
		if err != nil {
			return fmt.Errorf("cannot encrypt metadata: %v", err)
		}
		user.Metadata = encMetadata
	}
	if user.FirstName != "" {
		encFirstName, err := r.crypto.Encrypt(user.FirstName, encryptionKey)
		if err != nil {
			return fmt.Errorf("cannot encrypt first name: %v", err)
		}
		user.FirstName = encFirstName
	}
	if user.LastName != "" {
		encLastName, err := r.crypto.Encrypt(user.LastName, encryptionKey)
		if err != nil {
			return fmt.Errorf("cannot encrypt last name: %v", err)
		}
		user.LastName = encLastName
	}
	if user.Birthdate != "" {
		encBirthdate, err := r.crypto.Encrypt(user.Birthdate, encryptionKey)
		if err != nil {
			return fmt.Errorf("cannot encrypt birthdate: %v", err)
		}
		user.Birthdate = encBirthdate
	}
	user.EncryptedAt = time.Now().UTC()
	return nil
}
