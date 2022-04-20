package token_repository

import (
	"errors"
)

func (t *mongodbRepository) DecryptToken(token *Token) error {
	if token == nil {
		return errors.New("token is nil")
	}
	if len(t.internalEncryptionKeys) > 0 {
		encryptionKey, err := t.crypto.CombineSymmetricKeys(t.internalEncryptionKeys, token.InternalEncryptionLevel)
		if err != nil {
			return err
		}
		if token.Device != "" {
			decDevice, err := t.crypto.Decrypt(token.Device, encryptionKey)
			if err != nil {
				return err
			}
			token.Device = decDevice
		}
		if token.LoggedInFrom != "" {
			decLoggedInFrom, err := t.crypto.Decrypt(token.LoggedInFrom, encryptionKey)
			if err != nil {
				return err
			}
			token.LoggedInFrom = decLoggedInFrom
		}
	}
	return nil
}
