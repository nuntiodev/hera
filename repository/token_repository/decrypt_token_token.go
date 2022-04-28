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
		if token.DeviceInfo != "" {
			decDevice, err := t.crypto.Decrypt(token.DeviceInfo, encryptionKey)
			if err != nil {
				return err
			}
			token.DeviceInfo = decDevice
		}
		if token.LoggedInFrom.City != "" {
			decCity, err := t.crypto.Decrypt(token.LoggedInFrom.City, encryptionKey)
			if err != nil {
				return err
			}
			token.LoggedInFrom.City = decCity
		}
		if token.LoggedInFrom.Country != "" {
			decCountry, err := t.crypto.Decrypt(token.LoggedInFrom.Country, encryptionKey)
			if err != nil {
				return err
			}
			token.LoggedInFrom.Country = decCountry
		}
		if token.LoggedInFrom.CountryCode != "" {
			decCountryCode, err := t.crypto.Decrypt(token.LoggedInFrom.CountryCode, encryptionKey)
			if err != nil {
				return err
			}
			token.LoggedInFrom.CountryCode = decCountryCode
		}
	}
	return nil
}
