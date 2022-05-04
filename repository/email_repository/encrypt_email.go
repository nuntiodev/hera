package email_repository

import (
	"errors"
	"fmt"
)

func (t *defaultEmailRepository) EncryptEmail(action int, email *Email) error {
	if email == nil {
		return errors.New("encrypt: email is nil")
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
		encryptionKey, err = t.crypto.CombineSymmetricKeys(t.internalEncryptionKeys, int(email.InternalEncryptionLevel))
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("invalid action %d", action)
	}
	if email.Logo != "" {
		logoUrl, err := t.crypto.Encrypt(email.Logo, encryptionKey)
		if err != nil {
			return err
		}
		email.Logo = logoUrl
	}
	if email.WelcomeMessage != "" {
		welcomeMessage, err := t.crypto.Encrypt(email.WelcomeMessage, encryptionKey)
		if err != nil {
			return err
		}
		email.WelcomeMessage = welcomeMessage
	}
	if email.BodyMessage != "" {
		bodyMessage, err := t.crypto.Encrypt(email.BodyMessage, encryptionKey)
		if err != nil {
			return err
		}
		email.BodyMessage = bodyMessage
	}
	if email.FooterMessage != "" {
		footerMessage, err := t.crypto.Encrypt(email.FooterMessage, encryptionKey)
		if err != nil {
			return err
		}
		email.FooterMessage = footerMessage
	}
	if email.Subject != "" {
		subject, err := t.crypto.Encrypt(email.Subject, encryptionKey)
		if err != nil {
			return err
		}
		email.Subject = subject
	}
	if email.TemplatePath != "" {
		templatePath, err := t.crypto.Encrypt(email.TemplatePath, encryptionKey)
		if err != nil {
			return err
		}
		email.TemplatePath = templatePath
	}
	return nil
}
