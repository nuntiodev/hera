package user_mock

import (
	"errors"
	"fmt"
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
		return errors.New(fmt.Sprintf("different namespaces: user1: %s  user2: %s", userOne.Namespace, userTwo.Namespace))
	} else if userOne.Name != userTwo.Name {
		return errors.New(fmt.Sprintf("different names: user1: %s  user2: %s", userOne.Name, userTwo.Name))
	} else if userOne.Id != userTwo.Id {
		return errors.New(fmt.Sprintf("different ids: user1: %s  user2: %s", userOne.Id, userTwo.Id))
	} else if userOne.Email != userTwo.Email {
		return errors.New(fmt.Sprintf("different emails: user1: %s  user2: %s", userOne.Email, userTwo.Email))
	} else if (userOne.Birthdate != nil && userTwo.Birthdate != nil) && (userOne.Birthdate.Seconds != userTwo.Birthdate.Seconds) {
		return errors.New(fmt.Sprintf("different birthdays at: user1: %d  user2: %d", userOne.Birthdate.Seconds, userTwo.Birthdate.Seconds))
	} else if userOne.CreatedAt.Seconds != userTwo.CreatedAt.Seconds {
		return errors.New(fmt.Sprintf("different created at: user1: %d  user2: %d", userOne.CreatedAt.Seconds, userTwo.CreatedAt.Seconds))
	} else if userOne.Country != userTwo.Country {
		return errors.New(fmt.Sprintf("different countries: user1: %s  user2: %s", userOne.Country, userTwo.Country))
	} else if userOne.Image != userTwo.Image {
		return errors.New(fmt.Sprintf("different images: user1: %s  user2: %s", userOne.Image, userTwo.Image))
	} else if userOne.Gender != userTwo.Gender {
		return errors.New(fmt.Sprintf("different genders: user1: %s  user2: %s", userOne.Gender.String(), userTwo.Gender.String()))
	}
	return nil
}
