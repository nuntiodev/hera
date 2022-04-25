package config_repository

import "errors"

func (t *defaultConfigRepository) DecryptConfig(config *Config) error {
	if config == nil {
		return errors.New("config is nil")
	}
	if len(t.internalEncryptionKeys) > 0 {
		encryptionKey, err := t.crypto.CombineSymmetricKeys(t.internalEncryptionKeys, len(t.internalEncryptionKeys))
		if err != nil {
			return err
		}
		if config.Name != "" {
			decName, err := t.crypto.Decrypt(config.Name, encryptionKey)
			if err != nil {
				return err
			}
			config.Name = decName
		}
		if config.Website != "" {
			decWebsite, err := t.crypto.Decrypt(config.Website, encryptionKey)
			if err != nil {
				return err
			}
			config.Website = decWebsite
		}
		if config.About != "" {
			decAbout, err := t.crypto.Decrypt(config.About, encryptionKey)
			if err != nil {
				return err
			}
			config.About = decAbout
		}
		if config.Email != "" {
			decEmail, err := t.crypto.Decrypt(config.Email, encryptionKey)
			if err != nil {
				return err
			}
			config.Email = decEmail
		}
		if config.Logo != "" {
			decLogo, err := t.crypto.Encrypt(config.Logo, encryptionKey)
			if err != nil {
				return err
			}
			config.Logo = decLogo
		}
		if config.Terms != "" {
			decTerms, err := t.crypto.Decrypt(config.Terms, encryptionKey)
			if err != nil {
				return err
			}
			config.Terms = decTerms
		}
		if config.AuthConfig != nil {
			if config.AuthConfig.Logo != "" {
				decLogo, err := t.crypto.Decrypt(config.AuthConfig.Logo, encryptionKey)
				if err != nil {
					return err
				}
				config.AuthConfig.Logo = decLogo
			}
			if config.AuthConfig.WelcomeTitle != "" {
				decWelcomeTitle, err := t.crypto.Decrypt(config.AuthConfig.WelcomeTitle, encryptionKey)
				if err != nil {
					return err
				}
				config.AuthConfig.WelcomeTitle = decWelcomeTitle
			}
			if config.AuthConfig.WelcomeDetails != "" {
				decWelcomeDetails, err := t.crypto.Decrypt(config.AuthConfig.WelcomeDetails, encryptionKey)
				if err != nil {
					return err
				}
				config.AuthConfig.WelcomeDetails = decWelcomeDetails
			}
			if config.AuthConfig.LoginButton != "" {
				decLoginButton, err := t.crypto.Decrypt(config.AuthConfig.LoginButton, encryptionKey)
				if err != nil {
					return err
				}
				config.AuthConfig.LoginButton = decLoginButton
			}
			if config.AuthConfig.LoginTitle != "" {
				decLoginTitle, err := t.crypto.Decrypt(config.AuthConfig.LoginTitle, encryptionKey)
				if err != nil {
					return err
				}
				config.AuthConfig.LoginTitle = decLoginTitle
			}
			if config.AuthConfig.LoginDetails != "" {
				decLoginDetails, err := t.crypto.Decrypt(config.AuthConfig.LoginDetails, encryptionKey)
				if err != nil {
					return err
				}
				config.AuthConfig.LoginDetails = decLoginDetails
			}
			if config.AuthConfig.RegisterButton != "" {
				decRegisterButton, err := t.crypto.Decrypt(config.AuthConfig.RegisterButton, encryptionKey)
				if err != nil {
					return err
				}
				config.AuthConfig.RegisterButton = decRegisterButton
			}
			if config.AuthConfig.RegisterTitle != "" {
				decRegisterTitle, err := t.crypto.Decrypt(config.AuthConfig.RegisterTitle, encryptionKey)
				if err != nil {
					return err
				}
				config.AuthConfig.RegisterTitle = decRegisterTitle
			}
			if config.AuthConfig.RegisterDetails != "" {
				decRegisterDetails, err := t.crypto.Decrypt(config.AuthConfig.RegisterDetails, encryptionKey)
				if err != nil {
					return err
				}
				config.AuthConfig.RegisterDetails = decRegisterDetails
			}
			if config.AuthConfig.MissingPasswordTitle != "" {
				decMissingPasswordTitle, err := t.crypto.Decrypt(config.AuthConfig.MissingPasswordTitle, encryptionKey)
				if err != nil {
					return err
				}
				config.AuthConfig.MissingPasswordTitle = decMissingPasswordTitle
			}
			if config.AuthConfig.MissingPasswordDetails != "" {
				decMissingPasswordDetails, err := t.crypto.Decrypt(config.AuthConfig.MissingPasswordDetails, encryptionKey)
				if err != nil {
					return err
				}
				config.AuthConfig.MissingPasswordDetails = decMissingPasswordDetails
			}
			if config.AuthConfig.MissingEmailTitle != "" {
				decMissingEmailTitle, err := t.crypto.Decrypt(config.AuthConfig.MissingEmailTitle, encryptionKey)
				if err != nil {
					return err
				}
				config.AuthConfig.MissingEmailTitle = decMissingEmailTitle
			}
			if config.AuthConfig.MissingEmailDetails != "" {
				decMissingEmailDetails, err := t.crypto.Decrypt(config.AuthConfig.MissingEmailDetails, encryptionKey)
				if err != nil {
					return err
				}
				config.AuthConfig.MissingEmailDetails = decMissingEmailDetails
			}
			if config.AuthConfig.PasswordDoNotMatchTitle != "" {
				decPasswordsDoNotMatchTitle, err := t.crypto.Decrypt(config.AuthConfig.PasswordDoNotMatchTitle, encryptionKey)
				if err != nil {
					return err
				}
				config.AuthConfig.PasswordDoNotMatchTitle = decPasswordsDoNotMatchTitle
			}
			if config.AuthConfig.PasswordDoNotMatchDetails != "" {
				decPasswordsDoNotMatchDetails, err := t.crypto.Decrypt(config.AuthConfig.PasswordDoNotMatchDetails, encryptionKey)
				if err != nil {
					return err
				}
				config.AuthConfig.PasswordDoNotMatchDetails = decPasswordsDoNotMatchDetails
			}
		}
	}
	return nil
}
