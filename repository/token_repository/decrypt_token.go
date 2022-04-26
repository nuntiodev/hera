package token_repository

import (
	"errors"
	"fmt"
)

func (t *mongodbRepository) DecryptToken(token *Token) error {
	if token == nil {
		return errors.New("token is nil")
	}
	if len(t.internalEncryptionKeys) > 0 {
		encryptionKey, err := t.crypto.CombineSymmetricKeys(t.internalEncryptionKeys, token.InternalEncryptionLevel)
		if err != nil {
			fmt.Println("erjaisdjoi12 aw312")
			return err
		}
		if token.Device != "" {
			decDevice, err := t.crypto.Decrypt(token.Device, encryptionKey)
			if err != nil {
				fmt.Println("err is here 4")
				fmt.Println(token.Device)
				fmt.Println("err is here 4")
				return err
			}
			token.Device = decDevice
		}
		if token.LoggedInFrom.City != "" {
			decCity, err := t.crypto.Decrypt(token.LoggedInFrom.City, encryptionKey)
			if err != nil {
				fmt.Println("err is here 3")
				fmt.Println(token.LoggedInFrom.City)
				fmt.Println("err is here 3")
				return err
			}
			token.LoggedInFrom.City = decCity
		}
		if token.LoggedInFrom.Country != "" {
			decCountry, err := t.crypto.Decrypt(token.LoggedInFrom.Country, encryptionKey)
			if err != nil {
				fmt.Println("err is here 2")
				fmt.Println(token.LoggedInFrom.Country)
				fmt.Println("err is here 2")
				return err
			}
			token.LoggedInFrom.Country = decCountry
		}
		if token.LoggedInFrom.CountryCode != "" {
			decCountryCode, err := t.crypto.Decrypt(token.LoggedInFrom.CountryCode, encryptionKey)
			if err != nil {
				fmt.Println("err is here 1")
				fmt.Println(token.LoggedInFrom.CountryCode)
				fmt.Println("err is here 1")
				return err
			}
			token.LoggedInFrom.CountryCode = decCountryCode
		}
	}
	return nil
}
