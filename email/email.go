package email

import (
	"errors"
)

var (
	EmailSender Email
)

type Email interface {
	SendVerificationEmail(appName, to, code string) error
	SendResetPasswordEmail(appName, to, code string) error
}

func New() (Email, error) {
	if EmailSender == nil {
		return nil, errors.New("no email sender present")
	}
	return EmailSender, nil
}
