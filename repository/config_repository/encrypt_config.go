package config_repository

import (
	"errors"
)

func (t *defaultConfigRepository) EncryptConfig(action int, config *Config) error {
	if config == nil {
		return errors.New("config is nil")
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
		encryptionKey, err = t.crypto.CombineSymmetricKeys(t.internalEncryptionKeys, int(config.InternalEncryptionLevel))
		if err != nil {
			return err
		}
	default:
		return errors.New("invalid action")
	}
	if len(t.internalEncryptionKeys) > 0 {
		// set external encryption level
		config.InternalEncryptionLevel = int32(len(t.internalEncryptionKeys))
		if config.Name != "" {
			encName, err := t.crypto.Encrypt(config.Name, encryptionKey)
			if err != nil {
				return err
			}
			config.Name = encName
		}
		if config.Website != "" {
			encWebsite, err := t.crypto.Encrypt(config.Website, encryptionKey)
			if err != nil {
				return err
			}
			config.Website = encWebsite
		}
		if config.About != "" {
			encAbout, err := t.crypto.Encrypt(config.About, encryptionKey)
			if err != nil {
				return err
			}
			config.About = encAbout
		}
		if config.Email != "" {
			encEmail, err := t.crypto.Encrypt(config.Email, encryptionKey)
			if err != nil {
				return err
			}
			config.Email = encEmail
		}
		if config.Logo != "" {
			encLogo, err := t.crypto.Encrypt(config.Logo, encryptionKey)
			if err != nil {
				return err
			}
			config.Logo = encLogo
		}
		if config.Terms != "" {
			encTerms, err := t.crypto.Encrypt(config.Terms, encryptionKey)
			if err != nil {
				return err
			}
			config.Terms = encTerms
		}
		if config.AuthConfig != nil {
			if config.AuthConfig.WelcomeTitle != "" {
				encWelcomeTitle, err := t.crypto.Encrypt(config.AuthConfig.WelcomeTitle, encryptionKey)
				if err != nil {
					return err
				}
				config.AuthConfig.WelcomeTitle = encWelcomeTitle
			}
			if config.AuthConfig.WelcomeDetails != "" {
				encWelcomeDetails, err := t.crypto.Encrypt(config.AuthConfig.WelcomeDetails, encryptionKey)
				if err != nil {
					return err
				}
				config.AuthConfig.WelcomeDetails = encWelcomeDetails
			}
			if config.AuthConfig.LoginButton != "" {
				encLoginButton, err := t.crypto.Encrypt(config.AuthConfig.LoginButton, encryptionKey)
				if err != nil {
					return err
				}
				config.AuthConfig.LoginButton = encLoginButton
			}
			if config.AuthConfig.LoginTitle != "" {
				encLoginTitle, err := t.crypto.Encrypt(config.AuthConfig.LoginTitle, encryptionKey)
				if err != nil {
					return err
				}
				config.AuthConfig.LoginTitle = encLoginTitle
			}
			if config.AuthConfig.LoginDetails != "" {
				encLoginDetails, err := t.crypto.Encrypt(config.AuthConfig.LoginDetails, encryptionKey)
				if err != nil {
					return err
				}
				config.AuthConfig.LoginDetails = encLoginDetails
			}
			if config.AuthConfig.RegisterButton != "" {
				encRegisterButton, err := t.crypto.Encrypt(config.AuthConfig.RegisterButton, encryptionKey)
				if err != nil {
					return err
				}
				config.AuthConfig.RegisterButton = encRegisterButton
			}
			if config.AuthConfig.RegisterTitle != "" {
				encRegisterTitle, err := t.crypto.Encrypt(config.AuthConfig.RegisterTitle, encryptionKey)
				if err != nil {
					return err
				}
				config.AuthConfig.RegisterTitle = encRegisterTitle
			}
			if config.AuthConfig.RegisterDetails != "" {
				encRegisterDetails, err := t.crypto.Encrypt(config.AuthConfig.RegisterDetails, encryptionKey)
				if err != nil {
					return err
				}
				config.AuthConfig.RegisterDetails = encRegisterDetails
			}
			if config.AuthConfig.MissingPasswordTitle != "" {
				encMissingPasswordTitle, err := t.crypto.Encrypt(config.AuthConfig.MissingPasswordTitle, encryptionKey)
				if err != nil {
					return err
				}
				config.AuthConfig.MissingPasswordTitle = encMissingPasswordTitle
			}
			if config.AuthConfig.MissingPasswordDetails != "" {
				encMissingPasswordDetails, err := t.crypto.Encrypt(config.AuthConfig.MissingPasswordDetails, encryptionKey)
				if err != nil {
					return err
				}
				config.AuthConfig.MissingPasswordDetails = encMissingPasswordDetails
			}
			if config.AuthConfig.MissingEmailTitle != "" {
				encMissingEmailTitle, err := t.crypto.Encrypt(config.AuthConfig.MissingEmailTitle, encryptionKey)
				if err != nil {
					return err
				}
				config.AuthConfig.MissingEmailTitle = encMissingEmailTitle
			}
			if config.AuthConfig.MissingEmailDetails != "" {
				encMissingEmailDetails, err := t.crypto.Encrypt(config.AuthConfig.MissingEmailDetails, encryptionKey)
				if err != nil {
					return err
				}
				config.AuthConfig.MissingEmailDetails = encMissingEmailDetails
			}
			if config.AuthConfig.PasswordDoNotMatchTitle != "" {
				encPasswordsDoNotMatchTitle, err := t.crypto.Encrypt(config.AuthConfig.PasswordDoNotMatchTitle, encryptionKey)
				if err != nil {
					return err
				}
				config.AuthConfig.PasswordDoNotMatchTitle = encPasswordsDoNotMatchTitle
			}
			if config.AuthConfig.PasswordDoNotMatchDetails != "" {
				encPasswordsDoNotMatchDetails, err := t.crypto.Encrypt(config.AuthConfig.PasswordDoNotMatchDetails, encryptionKey)
				if err != nil {
					return err
				}
				config.AuthConfig.PasswordDoNotMatchDetails = encPasswordsDoNotMatchDetails
			}
		}
	}
	return nil
}
