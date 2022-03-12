package user_repository

import (
	"errors"
	"github.com/go-passwd/validator"
	hibp "github.com/mattevans/pwned-passwords"
	"github.com/softcorp-io/block-proto/go_block"
	ts "google.golang.org/protobuf/types/known/timestamppb"
)

var PwnedError = errors.New("this password has been involved in a data breach")

const (
	timeLayout = "2006-01-02 15:04:05"
)

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

func userToProtoUser(user *User) *go_block.User {
	if user == nil {
		return nil
	}
	return &go_block.User{
		Id:         user.Id,
		OptionalId: user.OptionalId,
		Namespace:  user.Namespace,
		Email:      user.Email,
		Role:       user.Role,
		Password:   user.Password,
		Image:      user.Image,
		Encrypted:  user.Encrypted,
		Metadata:   user.Metadata,
		CreatedAt:  ts.New(user.CreatedAt),
		UpdatedAt:  ts.New(user.UpdatedAt),
	}
}

func protoUserToUser(user *go_block.User) *User {
	if user == nil {
		return nil
	}
	return &User{
		Id:         user.Id,
		OptionalId: user.OptionalId,
		Namespace:  user.Namespace,
		Email:      user.Email,
		Role:       user.Role,
		Password:   user.Password,
		Image:      user.Image,
		Encrypted:  user.Encrypted,
		Metadata:   user.Metadata,
		CreatedAt:  user.CreatedAt.AsTime(),
		UpdatedAt:  user.UpdatedAt.AsTime(),
	}
}
