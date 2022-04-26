package token_repository

import "errors"

func (t *mongodbRepository) EncryptToken(action int, token *Token) error {
	if token == nil {
		return errors.New("token is nil")
	}
	encryptionKey := ""
	var err error
	switch action {
	case actionCreate:
		encryptionKey, err = t.crypto.CombineSymmetricKeys(t.internalEncryptionKeys, len(t.internalEncryptionKeys))
		if err != nil {
			return err
		}
	case actionUpdate:
		encryptionKey, err = t.crypto.CombineSymmetricKeys(t.internalEncryptionKeys, token.InternalEncryptionLevel)
		if err != nil {
			return err
		}
	}
	if len(t.internalEncryptionKeys) > 0 {
		if token.DeviceInfo != "" {
			encDevice, err := t.crypto.Encrypt(token.DeviceInfo, encryptionKey)
			if err != nil {
				return err
			}
			token.DeviceInfo = encDevice
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
	return nil
}
