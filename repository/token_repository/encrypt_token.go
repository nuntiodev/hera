package token_repository

import "errors"

func (t *mongodbRepository) EncryptToken(action int, token *Token) error {
	if token == nil {
		return errors.New("token is nil")
	}
	switch action {
	case actionCreate:
		if len(t.internalEncryptionKeys) > 0 {
			encryptionKey, err := t.crypto.CombineSymmetricKeys(t.internalEncryptionKeys, len(t.internalEncryptionKeys))
			if err != nil {
				return err
			}
			if token.Device != "" {
				encDevice, err := t.crypto.Encrypt(token.Device, encryptionKey)
				if err != nil {
					return err
				}
				token.Device = encDevice
			}
			if token.LoggedInFrom != "" {
				encLoggedInFrom, err := t.crypto.Encrypt(token.LoggedInFrom, encryptionKey)
				if err != nil {
					return err
				}
				token.LoggedInFrom = encLoggedInFrom
			}
		}
	}
	return nil
}
