package text_repository

import (
	"errors"
)

func (t *defaultTextRepository) DecryptText(text *Text) error {
	if text == nil {
		return errors.New("decrypt: text is nil")
	}
	if len(t.internalEncryptionKeys) > 0 {
		encryptionKey, err := t.crypto.CombineSymmetricKeys(t.internalEncryptionKeys, int(text.InternalEncryptionLevel))
		if err != nil {
			return err
		}
		if text.GeneralText != nil {

		}
		if text.GeneralText != nil {
			if text.GeneralText.MissingPasswordTitle != "" {
				missingPasswordTitle, err := t.crypto.Decrypt(text.GeneralText.MissingPasswordTitle, encryptionKey)
				if err != nil {
					return err
				}
				text.GeneralText.MissingPasswordTitle = missingPasswordTitle
			}
			if text.GeneralText.MissingPasswordDetails != "" {
				missingPasswordDetails, err := t.crypto.Decrypt(text.GeneralText.MissingPasswordDetails, encryptionKey)
				if err != nil {
					return err
				}
				text.GeneralText.MissingPasswordDetails = missingPasswordDetails
			}
			if text.GeneralText.MissingEmailTitle != "" {
				missingEmailTitle, err := t.crypto.Decrypt(text.GeneralText.MissingEmailTitle, encryptionKey)
				if err != nil {
					return err
				}
				text.GeneralText.MissingEmailTitle = missingEmailTitle
			}
			if text.GeneralText.MissingEmailDetails != "" {
				missingEmailDetails, err := t.crypto.Decrypt(text.GeneralText.MissingEmailDetails, encryptionKey)
				if err != nil {
					return err
				}
				text.GeneralText.MissingEmailDetails = missingEmailDetails
			}
			if text.GeneralText.CreatedBy != "" {
				createdBy, err := t.crypto.Decrypt(text.GeneralText.CreatedBy, encryptionKey)
				if err != nil {
					return err
				}
				text.GeneralText.CreatedBy = createdBy
			}
			if text.GeneralText.PasswordHint != "" {
				passwordHint, err := t.crypto.Decrypt(text.GeneralText.PasswordHint, encryptionKey)
				if err != nil {
					return err
				}
				text.GeneralText.PasswordHint = passwordHint
			}
			if text.GeneralText.EmailHint != "" {
				emailHint, err := t.crypto.Decrypt(text.GeneralText.EmailHint, encryptionKey)
				if err != nil {
					return err
				}
				text.GeneralText.EmailHint = emailHint
			}
			if text.GeneralText.ErrorTitle != "" {
				errorTitle, err := t.crypto.Decrypt(text.GeneralText.ErrorTitle, encryptionKey)
				if err != nil {
					return err
				}
				text.GeneralText.ErrorTitle = errorTitle
			}
			if text.GeneralText.ErrorDescription != "" {
				errorDescription, err := t.crypto.Decrypt(text.GeneralText.ErrorDescription, encryptionKey)
				if err != nil {
					return err
				}
				text.GeneralText.ErrorDescription = errorDescription
			}
			if text.GeneralText.NoWifiTitle != "" {
				noWifiTitle, err := t.crypto.Decrypt(text.GeneralText.NoWifiTitle, encryptionKey)
				if err != nil {
					return err
				}
				text.GeneralText.NoWifiTitle = noWifiTitle
			}
			if text.GeneralText.NoWifiDescription != "" {
				noWifiDescription, err := t.crypto.Decrypt(text.GeneralText.NoWifiDescription, encryptionKey)
				if err != nil {
					return err
				}
				text.GeneralText.NoWifiDescription = noWifiDescription
			}
		}
		if text.WelcomeText != nil {
			if text.WelcomeText.WelcomeTitle != "" {
				welcomeTitle, err := t.crypto.Decrypt(text.WelcomeText.WelcomeTitle, encryptionKey)
				if err != nil {
					return err
				}
				text.WelcomeText.WelcomeTitle = welcomeTitle
			}
			if text.WelcomeText.WelcomeDetails != "" {
				welcomeDetails, err := t.crypto.Decrypt(text.WelcomeText.WelcomeDetails, encryptionKey)
				if err != nil {
					return err
				}
				text.WelcomeText.WelcomeDetails = welcomeDetails
			}
			if text.WelcomeText.ContinueWithNuntio != "" {
				continueWithNuntio, err := t.crypto.Decrypt(text.WelcomeText.ContinueWithNuntio, encryptionKey)
				if err != nil {
					return err
				}
				text.WelcomeText.ContinueWithNuntio = continueWithNuntio
			}
		}
		if text.RegisterText != nil {
			if text.RegisterText.RegisterButton != "" {
				registerButton, err := t.crypto.Decrypt(text.RegisterText.RegisterButton, encryptionKey)
				if err != nil {
					return err
				}
				text.RegisterText.RegisterButton = registerButton
			}
			if text.RegisterText.RegisterTitle != "" {
				registerTitle, err := t.crypto.Decrypt(text.RegisterText.RegisterTitle, encryptionKey)
				if err != nil {
					return err
				}
				text.RegisterText.RegisterTitle = registerTitle
			}
			if text.RegisterText.RegisterDetails != "" {
				registerDetails, err := t.crypto.Decrypt(text.RegisterText.RegisterDetails, encryptionKey)
				if err != nil {
					return err
				}
				text.RegisterText.RegisterDetails = registerDetails
			}
			if text.RegisterText.PasswordDoNotMatchTitle != "" {
				passwordDoNotMatchTitle, err := t.crypto.Decrypt(text.RegisterText.PasswordDoNotMatchTitle, encryptionKey)
				if err != nil {
					return err
				}
				text.RegisterText.PasswordDoNotMatchTitle = passwordDoNotMatchTitle
			}
			if text.RegisterText.PasswordDoNotMatchDetails != "" {
				passwordDoNotMatchDetails, err := t.crypto.Decrypt(text.RegisterText.PasswordDoNotMatchDetails, encryptionKey)
				if err != nil {
					return err
				}
				text.RegisterText.PasswordDoNotMatchDetails = passwordDoNotMatchDetails
			}
			if text.RegisterText.RepeatPasswordHint != "" {
				repeatPasswordHint, err := t.crypto.Decrypt(text.RegisterText.RepeatPasswordHint, encryptionKey)
				if err != nil {
					return err
				}
				text.RegisterText.RepeatPasswordHint = repeatPasswordHint
			}
			if text.RegisterText.ContainsSpecialChar != "" {
				containsSpecialChar, err := t.crypto.Decrypt(text.RegisterText.ContainsSpecialChar, encryptionKey)
				if err != nil {
					return err
				}
				text.RegisterText.ContainsSpecialChar = containsSpecialChar
			}
			if text.RegisterText.ContainsNumberChar != "" {
				containsNumberChar, err := t.crypto.Decrypt(text.RegisterText.ContainsNumberChar, encryptionKey)
				if err != nil {
					return err
				}
				text.RegisterText.ContainsNumberChar = containsNumberChar
			}
			if text.RegisterText.PasswordMustMatch != "" {
				passwordMustMatch, err := t.crypto.Decrypt(text.RegisterText.PasswordMustMatch, encryptionKey)
				if err != nil {
					return err
				}
				text.RegisterText.PasswordMustMatch = passwordMustMatch
			}
			if text.RegisterText.ContainsEightChars != "" {
				containsEightChars, err := t.crypto.Decrypt(text.RegisterText.ContainsEightChars, encryptionKey)
				if err != nil {
					return err
				}
				text.RegisterText.ContainsEightChars = containsEightChars
			}
		}
		if text.LoginText != nil {
			if text.LoginText.LoginButton != "" {
				loginButton, err := t.crypto.Decrypt(text.LoginText.LoginButton, encryptionKey)
				if err != nil {
					return err
				}
				text.LoginText.LoginButton = loginButton
			}
			if text.LoginText.LoginTitle != "" {
				loginTitle, err := t.crypto.Decrypt(text.LoginText.LoginTitle, encryptionKey)
				if err != nil {
					return err
				}
				text.LoginText.LoginTitle = loginTitle
			}
			if text.LoginText.LoginDetails != "" {
				loginDetails, err := t.crypto.Decrypt(text.LoginText.LoginDetails, encryptionKey)
				if err != nil {
					return err
				}
				text.LoginText.LoginDetails = loginDetails
			}
			if text.LoginText.ForgotPassword != "" {
				forgotPassword, err := t.crypto.Decrypt(text.LoginText.ForgotPassword, encryptionKey)
				if err != nil {
					return err
				}
				text.LoginText.ForgotPassword = forgotPassword
			}
		}
		if text.ProfileText != nil {
			if text.ProfileText.ProfileTitle != "" {
				profileTitle, err := t.crypto.Decrypt(text.ProfileText.ProfileTitle, encryptionKey)
				if err != nil {
					return err
				}
				text.ProfileText.ProfileTitle = profileTitle
			}
			if text.ProfileText.Logout != "" {
				logout, err := t.crypto.Decrypt(text.ProfileText.Logout, encryptionKey)
				if err != nil {
					return err
				}
				text.ProfileText.Logout = logout
			}
			if text.ProfileText.ChangeEmailTitle != "" {
				changeEmailTitle, err := t.crypto.Decrypt(text.ProfileText.ChangeEmailTitle, encryptionKey)
				if err != nil {
					return err
				}
				text.ProfileText.ChangeEmailTitle = changeEmailTitle
			}
			if text.ProfileText.ChangeEmailDescription != "" {
				changeEmailDescription, err := t.crypto.Decrypt(text.ProfileText.ChangeEmailDescription, encryptionKey)
				if err != nil {
					return err
				}
				text.ProfileText.ChangeEmailDescription = changeEmailDescription
			}
			if text.ProfileText.ChangePasswordTitle != "" {
				changePasswordTitle, err := t.crypto.Decrypt(text.ProfileText.ChangePasswordTitle, encryptionKey)
				if err != nil {
					return err
				}
				text.ProfileText.ChangePasswordTitle = changePasswordTitle
			}
			if text.ProfileText.ChangePasswordDescription != "" {
				changePasswordDescription, err := t.crypto.Decrypt(text.ProfileText.ChangePasswordDescription, encryptionKey)
				if err != nil {
					return err
				}
				text.ProfileText.ChangePasswordDescription = changePasswordDescription
			}
		}
	}
	return nil
}
