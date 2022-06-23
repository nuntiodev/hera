package email

import (
	"errors"
)

var (
	EmailSender Email
)

type Email interface {
	SendVerificationEmail(to, code string) error
	SendResetPasswordEmail(to, code string) error
}

func New() (Email, error) {
	if EmailSender == nil {
		return nil, errors.New("no email sender present")
	}
	return EmailSender, nil
}
