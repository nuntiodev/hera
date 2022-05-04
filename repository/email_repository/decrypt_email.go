package email_repository

import (
	"errors"
)

func (t *defaultEmailRepository) DecryptEmail(email *Email) error {
	if email == nil {
		return errors.New("encrypt: email is nil")
	}
	encryptionKey, err := t.crypto.CombineSymmetricKeys(t.internalEncryptionKeys, int(email.InternalEncryptionLevel))
	if err != nil {
		return err
	}
	if email.Logo != "" {
		logoUrl, err := t.crypto.Decrypt(email.Logo, encryptionKey)
		if err != nil {
			return err
		}
		email.Logo = logoUrl
	}
	if email.WelcomeMessage != "" {
		welcomeMessage, err := t.crypto.Decrypt(email.WelcomeMessage, encryptionKey)
		if err != nil {
			return err
		}
		email.WelcomeMessage = welcomeMessage
	}
	if email.BodyMessage != "" {
		bodyMessage, err := t.crypto.Decrypt(email.BodyMessage, encryptionKey)
		if err != nil {
			return err
		}
		email.BodyMessage = bodyMessage
	}
	if email.FooterMessage != "" {
		footerMessage, err := t.crypto.Decrypt(email.FooterMessage, encryptionKey)
		if err != nil {
			return err
		}
		email.FooterMessage = footerMessage
	}
	if email.Subject != "" {
		subject, err := t.crypto.Decrypt(email.Subject, encryptionKey)
		if err != nil {
			return err
		}
		email.Subject = subject
	}
	if email.TemplatePath != "" {
		templatePath, err := t.crypto.Decrypt(email.TemplatePath, encryptionKey)
		if err != nil {
			return err
		}
		email.TemplatePath = templatePath
	}
	return nil
}
