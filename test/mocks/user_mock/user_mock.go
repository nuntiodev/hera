package user_mock

import (
	"errors"
	"github.com/brianvoe/gofakeit/v6"
	uuid "github.com/satori/go.uuid"
	"github.com/softcorp-io/block-proto/go_block/block_user"
	"math/rand"
)

func GetRandomGender() block_user.Gender {
	min := 0
	max := 2
	choice := rand.Intn(max-min) + min
	switch choice {
	case 0:
		return block_user.Gender_MALE
	case 1:
		return block_user.Gender_FEMALE
	case 2:
		return block_user.Gender_OTHER
	}
	return block_user.Gender_INVALID_GENDER
}

func GetRandomUser(user *block_user.User) *block_user.User {
	resp := &block_user.User{
		Email:     gofakeit.Email(),
		Password:  gofakeit.Password(true, true, true, true, true, 10),
		Namespace: uuid.NewV4().String(),
	}
	if user != nil {
		if user.Name != "" {
			resp.Name = user.Name
		}
		if user.Email != "" {
			resp.Email = user.Email
		}
		if user.Id != "" {
			resp.Id = user.Id
		}
		if user.Password != "" {
			resp.Password = user.Password
		}
		if user.Gender != block_user.Gender_INVALID_GENDER {
			resp.Gender = user.Gender
		}
		if user.Image != "" {
			resp.Image = user.Image
		}
		if user.Namespace != "" {
			resp.Namespace = user.Namespace
		}
		if user.Country != "" {
			resp.Country = user.Country
		}
		if user.Birthdate.IsValid() {
			resp.Birthdate = user.Birthdate
		}
		if user.Gender != block_user.Gender_INVALID_GENDER {
			resp.Gender = user.Gender
		}
	}
	return resp
}

func CompareUsers(userOne, userTwo *block_user.User) error {
	if userOne == nil || userTwo == nil {
		return errors.New("one of the users are nil")
	}
	if userOne.Namespace != userTwo.Namespace {
		return errors.New("different namespaces")
	} else if userOne.Name != userTwo.Name {
		return errors.New("different names")
	} else if userOne.Id != userTwo.Id {
		return errors.New("different ids")
	} else if userOne.Email != userTwo.Email {
		return errors.New("different emails")
	} else if userOne.Birthdate.String() != userTwo.Birthdate.String() {
		return errors.New("different birthdays")
	} else if userOne.CreatedAt.String() != userTwo.CreatedAt.String() {
		return errors.New("different created at")
	} else if userOne.Country != userTwo.Country {
		return errors.New("different countries")
	} else if userOne.Image != userTwo.Image {
		return errors.New("different images")
	} else if userOne.Gender != userTwo.Gender {
		return errors.New("different genders")
	}
	return nil
}
