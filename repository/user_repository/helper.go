package user_repository

import (
	"errors"
	"github.com/go-passwd/validator"
	hibp "github.com/mattevans/pwned-passwords"
	"github.com/softcorp-io/block-proto/go_block/block_user"
	ts "google.golang.org/protobuf/types/known/timestamppb"
	"time"
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

func genderToString(gender block_user.Gender) string {
	return gender.String()
}

func stringToGender(gender string) block_user.Gender {
	switch gender {
	case block_user.Gender_MALE.String():
		return block_user.Gender_MALE
	case block_user.Gender_FEMALE.String():
		return block_user.Gender_FEMALE
	case block_user.Gender_OTHER.String():
		return block_user.Gender_OTHER
	default:
		return block_user.Gender_INVALID_GENDER
	}
}

func stringToTime(t string) (time.Time, error) {
	t2, err := time.Parse(timeLayout, t)
	if err != nil {
		return time.Time{}, err
	}
	return t2, err
}

func timeToString(t time.Time) string {
	return t.Format(timeLayout)
}

func userToProtoUser(user *User) *block_user.User {
	birthdate, _ := stringToTime(user.Birthdate)
	return &block_user.User{
		Id:                        user.Id,
		OptionalId:                user.OptionalId,
		Namespace:                 user.Namespace,
		Role:                      user.Role,
		Name:                      user.Name,
		Email:                     user.Email,
		Password:                  user.Password,
		Gender:                    stringToGender(user.Gender),
		Country:                   user.Country,
		Image:                     user.Image,
		Blocked:                   user.Blocked,
		Verified:                  user.Verified,
		Encrypted:                 user.Encrypted,
		DisablePasswordValidation: user.DisablePasswordValidation,
		Birthdate:                 ts.New(birthdate),
		CreatedAt:                 ts.New(user.CreatedAt),
		UpdatedAt:                 ts.New(user.UpdatedAt),
	}
}

func protoUserToUser(user *block_user.User) *User {
	return &User{
		Id:                        user.Id,
		OptionalId:                user.OptionalId,
		Namespace:                 user.Namespace,
		Role:                      user.Role,
		Name:                      user.Name,
		Email:                     user.Email,
		Password:                  user.Password,
		Gender:                    genderToString(user.Gender),
		Country:                   user.Country,
		Image:                     user.Image,
		Blocked:                   user.Blocked,
		Verified:                  user.Verified,
		Encrypted:                 user.Encrypted,
		DisablePasswordValidation: user.DisablePasswordValidation,
		Birthdate:                 timeToString(user.Birthdate.AsTime()),
		CreatedAt:                 user.CreatedAt.AsTime(),
		UpdatedAt:                 user.UpdatedAt.AsTime(),
	}
}
