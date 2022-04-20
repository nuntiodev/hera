package token_repository

import (
	"errors"
	"fmt"
)

func (t *mongodbRepository) DecryptToken(token *Token) error {
	if token == nil {
		return errors.New("token is nil")
	}
	fmt.Println("get in hrekop12k3")
	if len(t.internalEncryptionKeys) > 0 {
		fmt.Println("get in asjdio123")
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
			decLocation, err := t.crypto.Decrypt(token.LoggedInFrom, encryptionKey)
			if err != nil {
				return err
			}
			token.Device = decLocation
		}
	}
	return nil
}
