package user_mock

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	uuid "github.com/satori/go.uuid"
	"github.com/softcorp-io/block-proto/go_block/block_user"
	"time"
)

type MetadataMock struct {
	Name      string
	Birthdate time.Time
	Gender    string
}

func GetMetadata(metadata *MetadataMock) string {
	meta := MetadataMock{
		Name:      gofakeit.Name(),
		Birthdate: time.Now(),
		Gender:    gofakeit.Gender(),
	}
	if metadata != nil {
		if metadata.Name != "" {
			meta.Name = metadata.Name
		}
		if metadata.Gender != "" {
			meta.Gender = metadata.Gender
		}
		if !metadata.Birthdate.IsZero() {
			meta.Birthdate = metadata.Birthdate
		}
	}
	metaString, err := json.Marshal(&meta)
	if err != nil {
		panic(err)
	}
	return string(metaString)
}

func GetRandomUser(user *block_user.User) *block_user.User {
	resp := &block_user.User{
		Password:  gofakeit.Password(true, true, true, true, true, 10),
		Namespace: uuid.NewV4().String(),
		Metadata:  GetMetadata(nil),
	}
	if user != nil {
		if user.Email != "" {
			resp.Email = user.Email
		}
		if user.Id != "" {
			resp.Id = user.Id
		}
		if user.OptionalId != "" {
			resp.OptionalId = user.OptionalId
		}
		if user.Role != "" {
			resp.Role = user.Role
		}
		if user.Password != "" {
			resp.Password = user.Password
		}
		if user.Image != "" {
			resp.Image = user.Image
		}
		if user.Namespace != "" {
			resp.Namespace = user.Namespace
		}
		if user.Metadata != "" {
			resp.Metadata = user.Metadata
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
	} else if userOne.Id != userTwo.Id {
		return errors.New(fmt.Sprintf("different ids: user1: %s  user2: %s", userOne.Id, userTwo.Id))
	} else if userOne.Email != userTwo.Email {
		return errors.New(fmt.Sprintf("different emails: user1: %s  user2: %s", userOne.Email, userTwo.Email))
	} else if userOne.CreatedAt.Seconds != userTwo.CreatedAt.Seconds {
		return errors.New(fmt.Sprintf("different created at: user1: %d  user2: %d", userOne.CreatedAt.Seconds, userTwo.CreatedAt.Seconds))
	} else if userOne.Image != userTwo.Image {
		return errors.New(fmt.Sprintf("different images: user1: %s  user2: %s", userOne.Image, userTwo.Image))
	} else if userOne.Encrypted != userTwo.Encrypted {
		return errors.New(fmt.Sprintf("different encrypted: user1: %t  user2: %t", userOne.Encrypted, userTwo.Encrypted))
	}
	if userOne.Encrypted == false {
		var metadataOne MetadataMock
		if err := json.Unmarshal([]byte(userOne.Metadata), &metadataOne); err != nil {
			return err
		}
		var metadataTwo MetadataMock
		if err := json.Unmarshal([]byte(userTwo.Metadata), &metadataTwo); err != nil {
			return err
		}
		if metadataOne.Name != metadataTwo.Name {
			return errors.New(fmt.Sprintf("different metadata names: user1: %s  user2: %s", metadataOne.Name, metadataTwo.Name))
		}
		if metadataOne.Gender != metadataTwo.Gender {
			return errors.New(fmt.Sprintf("different metadata genders: user1: %s  user2: %s", metadataOne.Gender, metadataTwo.Gender))
		}
		if metadataOne.Birthdate != metadataTwo.Birthdate {
			return errors.New(fmt.Sprintf("different metadata genders: user1: %s  user2: %s", metadataOne.Birthdate.String(), metadataTwo.Birthdate.String()))
		}
	}
	return nil
}
