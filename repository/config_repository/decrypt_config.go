package config_repository

import "errors"

func (t *defaultConfigRepository) DecryptConfig(config *Config) error {
	if config == nil {
		return errors.New("decrypt: config is nil")
	}
	if len(t.internalEncryptionKeys) > 0 {
		encryptionKey, err := t.crypto.CombineSymmetricKeys(t.internalEncryptionKeys, len(t.internalEncryptionKeys))
		if err != nil {
			return err
		}
		if config.Name != "" {
			name, err := t.crypto.Decrypt(config.Name, encryptionKey)
			if err != nil {
				return err
			}
			config.Name = name
		}
		if config.Website != "" {
			website, err := t.crypto.Decrypt(config.Website, encryptionKey)
			if err != nil {
				return err
			}
			config.Website = website
		}
		if config.About != "" {
			about, err := t.crypto.Decrypt(config.About, encryptionKey)
			if err != nil {
				return err
			}
			config.About = about
		}
		if config.Email != "" {
			email, err := t.crypto.Decrypt(config.Email, encryptionKey)
			if err != nil {
				return err
			}
			config.Email = email
		}
		if config.Logo != "" {
			logo, err := t.crypto.Encrypt(config.Logo, encryptionKey)
			if err != nil {
				return err
			}
			config.Logo = logo
		}
		if config.Terms != "" {
			terms, err := t.crypto.Decrypt(config.Terms, encryptionKey)
			if err != nil {
				return err
			}
			config.Terms = terms
		}
		if config.GeneralText != nil {
			if config.GeneralText.MissingPasswordTitle != "" {
				missingPasswordTitle, err := t.crypto.Decrypt(config.GeneralText.MissingPasswordTitle, encryptionKey)
				if err != nil {
					return err
				}
				config.GeneralText.MissingPasswordTitle = missingPasswordTitle
			}
			if config.GeneralText.MissingPasswordDetails != "" {
				missingPasswordDetails, err := t.crypto.Decrypt(config.GeneralText.MissingPasswordDetails, encryptionKey)
				if err != nil {
					return err
				}
				config.GeneralText.MissingPasswordDetails = missingPasswordDetails
			}
			if config.GeneralText.MissingEmailTitle != "" {
				missingEmailTitle, err := t.crypto.Decrypt(config.GeneralText.MissingEmailTitle, encryptionKey)
				if err != nil {
					return err
				}
				config.GeneralText.MissingEmailTitle = missingEmailTitle
			}
			if config.GeneralText.MissingEmailDetails != "" {
				missingEmailDetails, err := t.crypto.Decrypt(config.GeneralText.MissingEmailDetails, encryptionKey)
				if err != nil {
					return err
				}
				config.GeneralText.MissingEmailDetails = missingEmailDetails
			}
			if config.GeneralText.CreatedBy != "" {
				createdBy, err := t.crypto.Decrypt(config.GeneralText.CreatedBy, encryptionKey)
				if err != nil {
					return err
				}
				config.GeneralText.CreatedBy = createdBy
			}
			if config.GeneralText.PasswordHint != "" {
				passwordHint, err := t.crypto.Decrypt(config.GeneralText.PasswordHint, encryptionKey)
				if err != nil {
					return err
				}
				config.GeneralText.PasswordHint = passwordHint
			}
			if config.GeneralText.EmailHint != "" {
				emailHint, err := t.crypto.Decrypt(config.GeneralText.EmailHint, encryptionKey)
				if err != nil {
					return err
				}
				config.GeneralText.EmailHint = emailHint
			}
			if config.GeneralText.ErrorTitle != "" {
				errorTitle, err := t.crypto.Decrypt(config.GeneralText.ErrorTitle, encryptionKey)
				if err != nil {
					return err
				}
				config.GeneralText.ErrorTitle = errorTitle
			}
			if config.GeneralText.ErrorDescription != "" {
				errorDescription, err := t.crypto.Decrypt(config.GeneralText.ErrorDescription, encryptionKey)
				if err != nil {
					return err
				}
				config.GeneralText.ErrorDescription = errorDescription
			}
			if config.GeneralText.NoWifiTitle != "" {
				noWifiTitle, err := t.crypto.Decrypt(config.GeneralText.NoWifiTitle, encryptionKey)
				if err != nil {
					return err
				}
				config.GeneralText.NoWifiTitle = noWifiTitle
			}
			if config.GeneralText.NoWifiDescription != "" {
				noWifiDescription, err := t.crypto.Decrypt(config.GeneralText.NoWifiDescription, encryptionKey)
				if err != nil {
					return err
				}
				config.GeneralText.NoWifiDescription = noWifiDescription
			}
		}
		if config.WelcomeText != nil {
			if config.WelcomeText.WelcomeTitle != "" {
				welcomeTitle, err := t.crypto.Decrypt(config.WelcomeText.WelcomeTitle, encryptionKey)
				if err != nil {
					return err
				}
				config.WelcomeText.WelcomeTitle = welcomeTitle
			}
			if config.WelcomeText.WelcomeDetails != "" {
				welcomeDetails, err := t.crypto.Decrypt(config.WelcomeText.WelcomeDetails, encryptionKey)
				if err != nil {
					return err
				}
				config.WelcomeText.WelcomeDetails = welcomeDetails
			}
		}
		if config.RegisterText != nil {
			if config.RegisterText.RegisterButton != "" {
				registerButton, err := t.crypto.Decrypt(config.RegisterText.RegisterButton, encryptionKey)
				if err != nil {
					return err
				}
				config.RegisterText.RegisterButton = registerButton
			}
			if config.RegisterText.RegisterTitle != "" {
				registerTitle, err := t.crypto.Decrypt(config.RegisterText.RegisterTitle, encryptionKey)
				if err != nil {
					return err
				}
				config.RegisterText.RegisterTitle = registerTitle
			}
			if config.RegisterText.RegisterDetails != "" {
				registerDetails, err := t.crypto.Decrypt(config.RegisterText.RegisterDetails, encryptionKey)
				if err != nil {
					return err
				}
				config.RegisterText.RegisterDetails = registerDetails
			}
			if config.RegisterText.PasswordDoNotMatchTitle != "" {
				passwordDoNotMatchTitle, err := t.crypto.Decrypt(config.RegisterText.PasswordDoNotMatchTitle, encryptionKey)
				if err != nil {
					return err
				}
				config.RegisterText.PasswordDoNotMatchTitle = passwordDoNotMatchTitle
			}
			if config.RegisterText.PasswordDoNotMatchDetails != "" {
				passwordDoNotMatchDetails, err := t.crypto.Decrypt(config.RegisterText.PasswordDoNotMatchDetails, encryptionKey)
				if err != nil {
					return err
				}
				config.RegisterText.PasswordDoNotMatchDetails = passwordDoNotMatchDetails
			}
			if config.RegisterText.RepeatPasswordHint != "" {
				repeatPasswordHint, err := t.crypto.Decrypt(config.RegisterText.RepeatPasswordHint, encryptionKey)
				if err != nil {
					return err
				}
				config.RegisterText.RepeatPasswordHint = repeatPasswordHint
			}
			if config.RegisterText.ContainsSpecialChar != "" {
				containsSpecialChar, err := t.crypto.Decrypt(config.RegisterText.ContainsSpecialChar, encryptionKey)
				if err != nil {
					return err
				}
				config.RegisterText.ContainsSpecialChar = containsSpecialChar
			}
			if config.RegisterText.ContainsNumberChar != "" {
				containsNumberChar, err := t.crypto.Decrypt(config.RegisterText.ContainsNumberChar, encryptionKey)
				if err != nil {
					return err
				}
				config.RegisterText.ContainsNumberChar = containsNumberChar
			}
			if config.RegisterText.PasswordMustMatch != "" {
				passwordMustMatch, err := t.crypto.Decrypt(config.RegisterText.PasswordMustMatch, encryptionKey)
				if err != nil {
					return err
				}
				config.RegisterText.PasswordMustMatch = passwordMustMatch
			}
			if config.RegisterText.ContainsEightChars != "" {
				containsEightChars, err := t.crypto.Decrypt(config.RegisterText.ContainsEightChars, encryptionKey)
				if err != nil {
					return err
				}
				config.RegisterText.ContainsEightChars = containsEightChars
			}
		}
		if config.LoginText != nil {
			if config.LoginText.LoginButton != "" {
				loginButton, err := t.crypto.Decrypt(config.LoginText.LoginButton, encryptionKey)
				if err != nil {
					return err
				}
				config.LoginText.LoginButton = loginButton
			}
			if config.LoginText.LoginTitle != "" {
				loginTitle, err := t.crypto.Decrypt(config.LoginText.LoginTitle, encryptionKey)
				if err != nil {
					return err
				}
				config.LoginText.LoginTitle = loginTitle
			}
			if config.LoginText.LoginDetails != "" {
				loginDetails, err := t.crypto.Decrypt(config.LoginText.LoginDetails, encryptionKey)
				if err != nil {
					return err
				}
				config.LoginText.LoginDetails = loginDetails
			}
			if config.LoginText.ForgotPassword != "" {
				forgotPassword, err := t.crypto.Decrypt(config.LoginText.ForgotPassword, encryptionKey)
				if err != nil {
					return err
				}
				config.LoginText.ForgotPassword = forgotPassword
			}
		}
	}
	return nil
}
