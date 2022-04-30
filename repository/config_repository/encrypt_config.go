package config_repository

import (
	"errors"
)

func (t *defaultConfigRepository) EncryptConfig(action int, config *Config) error {
	if config == nil {
		return errors.New("encrypt: config is nil")
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
	// set external encryption level
	config.InternalEncryptionLevel = int32(len(t.internalEncryptionKeys))
	if config.Name != "" {
		name, err := t.crypto.Encrypt(config.Name, encryptionKey)
		if err != nil {
			return err
		}
		config.Name = name
	}
	if config.Website != "" {
		website, err := t.crypto.Encrypt(config.Website, encryptionKey)
		if err != nil {
			return err
		}
		config.Website = website
	}
	if config.About != "" {
		about, err := t.crypto.Encrypt(config.About, encryptionKey)
		if err != nil {
			return err
		}
		config.About = about
	}
	if config.Email != "" {
		email, err := t.crypto.Encrypt(config.Email, encryptionKey)
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
		terms, err := t.crypto.Encrypt(config.Terms, encryptionKey)
		if err != nil {
			return err
		}
		config.Terms = terms
	}
	if config.GeneralText != nil {
		if config.GeneralText.MissingPasswordTitle != "" {
			missingPasswordTitle, err := t.crypto.Encrypt(config.GeneralText.MissingPasswordTitle, encryptionKey)
			if err != nil {
				return err
			}
			config.GeneralText.MissingPasswordTitle = missingPasswordTitle
		}
		if config.GeneralText.MissingPasswordDetails != "" {
			missingPasswordDetails, err := t.crypto.Encrypt(config.GeneralText.MissingPasswordDetails, encryptionKey)
			if err != nil {
				return err
			}
			config.GeneralText.MissingPasswordDetails = missingPasswordDetails
		}
		if config.GeneralText.MissingEmailTitle != "" {
			missingEmailTitle, err := t.crypto.Encrypt(config.GeneralText.MissingEmailTitle, encryptionKey)
			if err != nil {
				return err
			}
			config.GeneralText.MissingEmailTitle = missingEmailTitle
		}
		if config.GeneralText.MissingEmailDetails != "" {
			missingEmailDetails, err := t.crypto.Encrypt(config.GeneralText.MissingEmailDetails, encryptionKey)
			if err != nil {
				return err
			}
			config.GeneralText.MissingEmailDetails = missingEmailDetails
		}
		if config.GeneralText.CreatedBy != "" {
			createdBy, err := t.crypto.Encrypt(config.GeneralText.CreatedBy, encryptionKey)
			if err != nil {
				return err
			}
			config.GeneralText.CreatedBy = createdBy
		}
		if config.GeneralText.PasswordHint != "" {
			passwordHint, err := t.crypto.Encrypt(config.GeneralText.PasswordHint, encryptionKey)
			if err != nil {
				return err
			}
			config.GeneralText.PasswordHint = passwordHint
		}
		if config.GeneralText.EmailHint != "" {
			emailHint, err := t.crypto.Encrypt(config.GeneralText.EmailHint, encryptionKey)
			if err != nil {
				return err
			}
			config.GeneralText.EmailHint = emailHint
		}
		if config.GeneralText.ErrorTitle != "" {
			errorTitle, err := t.crypto.Encrypt(config.GeneralText.ErrorTitle, encryptionKey)
			if err != nil {
				return err
			}
			config.GeneralText.ErrorTitle = errorTitle
		}
		if config.GeneralText.ErrorDescription != "" {
			errorDescription, err := t.crypto.Encrypt(config.GeneralText.ErrorDescription, encryptionKey)
			if err != nil {
				return err
			}
			config.GeneralText.ErrorDescription = errorDescription
		}
		if config.GeneralText.NoWifiTitle != "" {
			noWifiTitle, err := t.crypto.Encrypt(config.GeneralText.NoWifiTitle, encryptionKey)
			if err != nil {
				return err
			}
			config.GeneralText.NoWifiTitle = noWifiTitle
		}
		if config.GeneralText.NoWifiDescription != "" {
			noWifiDescription, err := t.crypto.Encrypt(config.GeneralText.NoWifiDescription, encryptionKey)
			if err != nil {
				return err
			}
			config.GeneralText.NoWifiDescription = noWifiDescription
		}
	}
	if config.WelcomeText != nil {
		if config.WelcomeText.WelcomeTitle != "" {
			welcomeTitle, err := t.crypto.Encrypt(config.WelcomeText.WelcomeTitle, encryptionKey)
			if err != nil {
				return err
			}
			config.WelcomeText.WelcomeTitle = welcomeTitle
		}
		if config.WelcomeText.WelcomeDetails != "" {
			welcomeDetails, err := t.crypto.Encrypt(config.WelcomeText.WelcomeDetails, encryptionKey)
			if err != nil {
				return err
			}
			config.WelcomeText.WelcomeDetails = welcomeDetails
		}
	}
	if config.RegisterText != nil {
		if config.RegisterText.RegisterButton != "" {
			registerButton, err := t.crypto.Encrypt(config.RegisterText.RegisterButton, encryptionKey)
			if err != nil {
				return err
			}
			config.RegisterText.RegisterButton = registerButton
		}
		if config.RegisterText.RegisterTitle != "" {
			registerTitle, err := t.crypto.Encrypt(config.RegisterText.RegisterTitle, encryptionKey)
			if err != nil {
				return err
			}
			config.RegisterText.RegisterTitle = registerTitle
		}
		if config.RegisterText.RegisterDetails != "" {
			registerDetails, err := t.crypto.Encrypt(config.RegisterText.RegisterDetails, encryptionKey)
			if err != nil {
				return err
			}
			config.RegisterText.RegisterDetails = registerDetails
		}
		if config.RegisterText.PasswordDoNotMatchTitle != "" {
			passwordDoNotMatchTitle, err := t.crypto.Encrypt(config.RegisterText.PasswordDoNotMatchTitle, encryptionKey)
			if err != nil {
				return err
			}
			config.RegisterText.PasswordDoNotMatchTitle = passwordDoNotMatchTitle
		}
		if config.RegisterText.PasswordDoNotMatchDetails != "" {
			passwordDoNotMatchDetails, err := t.crypto.Encrypt(config.RegisterText.PasswordDoNotMatchDetails, encryptionKey)
			if err != nil {
				return err
			}
			config.RegisterText.PasswordDoNotMatchDetails = passwordDoNotMatchDetails
		}
		if config.RegisterText.RepeatPasswordHint != "" {
			repeatPasswordHint, err := t.crypto.Encrypt(config.RegisterText.RepeatPasswordHint, encryptionKey)
			if err != nil {
				return err
			}
			config.RegisterText.RepeatPasswordHint = repeatPasswordHint
		}
		if config.RegisterText.ContainsNumberChar != "" {
			containsNumberChar, err := t.crypto.Encrypt(config.RegisterText.ContainsNumberChar, encryptionKey)
			if err != nil {
				return err
			}
			config.RegisterText.ContainsNumberChar = containsNumberChar
		}
		if config.RegisterText.PasswordMustMatch != "" {
			passwordMustMatch, err := t.crypto.Encrypt(config.RegisterText.PasswordMustMatch, encryptionKey)
			if err != nil {
				return err
			}
			config.RegisterText.PasswordMustMatch = passwordMustMatch
		}
		if config.RegisterText.ContainsEightChars != "" {
			containsEightChars, err := t.crypto.Encrypt(config.RegisterText.ContainsEightChars, encryptionKey)
			if err != nil {
				return err
			}
			config.RegisterText.ContainsEightChars = containsEightChars
		}
	}
	if config.LoginText != nil {
		if config.LoginText.LoginButton != "" {
			loginButton, err := t.crypto.Encrypt(config.LoginText.LoginButton, encryptionKey)
			if err != nil {
				return err
			}
			config.LoginText.LoginButton = loginButton
		}
		if config.LoginText.LoginTitle != "" {
			loginTitle, err := t.crypto.Encrypt(config.LoginText.LoginTitle, encryptionKey)
			if err != nil {
				return err
			}
			config.LoginText.LoginTitle = loginTitle
		}
		if config.LoginText.LoginDetails != "" {
			loginDetails, err := t.crypto.Encrypt(config.LoginText.LoginDetails, encryptionKey)
			if err != nil {
				return err
			}
			config.LoginText.LoginDetails = loginDetails
		}
		if config.LoginText.ForgotPassword != "" {
			forgotPassword, err := t.crypto.Encrypt(config.LoginText.ForgotPassword, encryptionKey)
			if err != nil {
				return err
			}
			config.LoginText.ForgotPassword = forgotPassword
		}
	}
	return nil
}
