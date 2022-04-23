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
			if token.LoggedInFrom.City != "" {
				encCity, err := t.crypto.Encrypt(token.LoggedInFrom.City, encryptionKey)
				if err != nil {
					return err
				}
				token.LoggedInFrom.City = encCity
			}
			if token.LoggedInFrom.Country != "" {
				encCountry, err := t.crypto.Encrypt(token.LoggedInFrom.Country, encryptionKey)
				if err != nil {
					return err
				}
				token.LoggedInFrom.Country = encCountry
			}
			if token.LoggedInFrom.CountryCode != "" {
				encCountryCode, err := t.crypto.Encrypt(token.LoggedInFrom.CountryCode, encryptionKey)
				if err != nil {
					return err
				}
				token.LoggedInFrom.CountryCode = encCountryCode
			}
		}
	}
	return nil
}
