package user_repository

import (
	"errors"
	"github.com/go-passwd/validator"
	hibp "github.com/mattevans/pwned-passwords"
)

var PwnedError = errors.New("this password has been involved in a data breach")

func validatePassword(password string) error {
	passwordValidator := validator.New(
		validator.MinLength(10, errors.New("password needs to contain at least 5 chars")),
		validator.MaxLength(100, errors.New("password needs to contain at below 100 chars")),
		validator.ContainsAtLeast("0123456789", 1, errors.New("password needs to contain at least one number")),
	)
	if err := passwordValidator.Validate(password); err != nil {
		return err
	}
	client := hibp.NewClient()
	pwned, err := client.Compromised(password)
	if err != nil {
		return err
	}
	if pwned {
		return errors.New("this password has been involved in a data breach")
	}
	return nil
}
