package user_repository

import (
	"errors"
	"github.com/go-passwd/validator"
	hibp "github.com/mattevans/pwned-passwords"
	"github.com/softcorp-io/block-proto/go_block/block_user"
	ts "google.golang.org/protobuf/types/known/timestamppb"
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

func userToProtoUser(user *User) *block_user.User {
	return &block_user.User{
		Id:                    user.Id,
		OptionalId:            user.OptionalId,
		Namespace:             user.Namespace,
		Role:                  user.Role,
		Name:                  user.Name,
		Email:                 user.Email,
		Password:              user.Password,
		Gender:                user.Gender,
		Country:               user.Country,
		Image:                 user.Image,
		Blocked:               user.Blocked,
		Verified:              user.Verified,
		DisableAuthentication: user.DisableAuthentication,
		Birthdate:             ts.New(user.Birthdate),
		CreatedAt:             ts.New(user.CreatedAt),
		UpdatedAt:             ts.New(user.UpdatedAt),
	}
}

func protoUserToUser(user *block_user.User) *User {
	return &User{
		Id:                    user.Id,
		OptionalId:            user.OptionalId,
		Namespace:             user.Namespace,
		Role:                  user.Role,
		Name:                  user.Name,
		Email:                 user.Email,
		Password:              user.Password,
		Gender:                user.Gender,
		Country:               user.Country,
		Image:                 user.Image,
		Blocked:               user.Blocked,
		Verified:              user.Verified,
		DisableAuthentication: user.DisableAuthentication,
		Birthdate:             user.Birthdate.AsTime(),
		CreatedAt:             user.CreatedAt.AsTime(),
		UpdatedAt:             user.UpdatedAt.AsTime(),
	}
}
