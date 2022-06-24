package text

import (
	"errors"
)

var (
	TextSender Text
)

type Text interface {
	SendVerificationText(appName, to, code string) error
	SendResetPasswordText(appName, to, code string) error
}

func New() (Text, error) {
	if TextSender == nil {
		return nil, errors.New("no text sender present")
	}
	return TextSender, nil
}
