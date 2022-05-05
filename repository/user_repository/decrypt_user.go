package user_repository

import (
	"errors"
	"fmt"
)

func (r *mongodbRepository) decryptUser(user *User) error {
	if user == nil {
		return errors.New("user is nil")
	}
	// encrypt using internal keys first
	if user.InternalEncryptionLevel > 0 {
		if len(r.internalEncryptionKeys) >= user.InternalEncryptionLevel {
			internalKey, err := r.crypto.CombineSymmetricKeys(r.internalEncryptionKeys, user.InternalEncryptionLevel)
			if err != nil {
				return err
			}
			if err := r.decrypt(user, internalKey); err != nil {
				return fmt.Errorf("could not decrypt user using internal key: %v", err)
			}
		} else {
			return errors.New("missing required internal encryption keys")
		}
	}
	if user.ExternalEncryptionLevel > 0 && r.externalEncryptionKey != "" {
		if err := r.decrypt(user, r.externalEncryptionKey); err != nil {
			return fmt.Errorf("could not decrypt user using external key: %v", err)
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
			return fmt.Errorf("cannot decrypt email: %v", err)
		}
		user.Email = decEmail
	}
	if user.Image != "" {
		decImage, err := r.crypto.Decrypt(user.Image, encryptionKey)
		if err != nil {
			return fmt.Errorf("cannot decrypt image: %v", err)
			return err
		}
		user.Image = decImage
	}
	if user.Metadata != "" {
		decMetadata, err := r.crypto.Decrypt(user.Metadata, encryptionKey)
		if err != nil {
			return fmt.Errorf("cannot decrypt metadata: %v", err)
		}
		user.Metadata = decMetadata
	}
	if user.FirstName != "" {
		decFirstName, err := r.crypto.Decrypt(user.FirstName, encryptionKey)
		if err != nil {
			return fmt.Errorf("cannot decrypt first name: %v", err)
		}
		user.FirstName = decFirstName
	}
	if user.LastName != "" {
		decLastName, err := r.crypto.Decrypt(user.LastName, encryptionKey)
		if err != nil {
			return fmt.Errorf("cannot decrypt last name: %v", err)
		}
		user.LastName = decLastName
	}
	if user.Birthdate != "" {
		decBirthdate, err := r.crypto.Decrypt(user.Birthdate, encryptionKey)
		if err != nil {
			return fmt.Errorf("cannot decrypt birthdate: %v", err)
		}
		user.Birthdate = decBirthdate
	}
	return nil
}
