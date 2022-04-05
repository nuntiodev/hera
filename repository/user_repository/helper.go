package user_repository

import (
	"encoding/json"
	"errors"
	"github.com/badoux/checkmail"
	uuid "github.com/satori/go.uuid"
	"github.com/softcorp-io/block-proto/go_block"
	passwordvalidator "github.com/wagslane/go-password-validator"
	ts "google.golang.org/protobuf/types/known/timestamppb"
	"strings"
)

var PwnedError = errors.New("this password has been involved in a data breach")

func prepare(action int, user *go_block.User) {
	if user == nil {
		return
	}
	switch action {
	case actionCreate:
		user.CreatedAt = ts.Now()
		user.UpdatedAt = ts.Now()
		if strings.TrimSpace(user.Id) == "" {
			user.Id = uuid.NewV4().String()
		}
	case actionUpdatePassword, actionUpdateImage, actionUpdateMetadata,
		actionUpdateNamespace, actionUpdateSecurity, actionUpdateEmail,
		actionUpdateOptionalId:
		user.UpdatedAt = ts.Now()
	}
	user.Id = strings.TrimSpace(user.Id)
	user.Email = strings.TrimSpace(strings.ToLower(user.Email))
	user.Image = strings.TrimSpace(user.Image)
	user.OptionalId = strings.TrimSpace(user.OptionalId)
	user.Metadata = strings.TrimSpace(user.Metadata)
}

func (r *mongoRepository) validate(action int, user *go_block.User) error {
	if user == nil {
		return errors.New("user is nil")
	}
	switch action {
	case actionGet:
		if user.Id == "" && user.Email == "" && user.OptionalId == "" {
			return errors.New("missing required search parameter")
		}
	case actionCreate:
		if user.Id == "" {
			return errors.New("invalid user id")
		} else if err := checkmail.ValidateFormat(user.Email); user.Email != "" && err != nil {
			return err
		} else if err := validatePassword(user.Password); err != nil && r.validatePassword {
			return err
		} else if !user.CreatedAt.IsValid() {
			return errors.New("invalid created at date")
		} else if !user.UpdatedAt.IsValid() {
			return errors.New("invalid updated at date")
		} else if !json.Valid([]byte(user.Metadata)) && user.Metadata != "" {
			return errors.New("invalid json type")
		}
	case actionUpdatePassword:
		if err := validatePassword(user.Password); err != nil && r.validatePassword {
			return err
		} else if !user.UpdatedAt.IsValid() {
			return errors.New("invalid updated at")
		}
	case actionUpdateEmail:
		if err := checkmail.ValidateFormat(user.Email); user.Email != "" && err != nil {
			return err
		} else if !user.UpdatedAt.IsValid() {
			return errors.New("invalid updated at")
		}
	case actionUpdateMetadata:
		if !user.UpdatedAt.IsValid() {
			return errors.New("invalid updated at")
		} else if !json.Valid([]byte(user.Metadata)) && user.Metadata != "" {
			return errors.New("invalid json type")
		}
	case actionUpdateSecurity:
		if !user.UpdatedAt.IsValid() {
			return errors.New("invalid updated at")
		}
	case actionGetAll, actionUpdateOptionalId:
		return nil
	}
	if len(user.Email) > maxFieldLength {

	} else if len(user.OptionalId) > maxFieldLength {

	} else if len(user.Metadata) > 10*maxFieldLength {

	}
	return nil
}

func validatePassword(password string) error {
	if err := passwordvalidator.Validate(password, minEntropy); err != nil {
		return err
	}
	return nil
}

func UserToProtoUser(user *User) *go_block.User {
	if user == nil {
		return nil
	}
	return &go_block.User{
		Id:                      user.Id,
		OptionalId:              user.OptionalId,
		Email:                   user.Email,
		Password:                user.Password,
		Image:                   user.Image,
		ExternalEncrypted:       user.ExternalEncrypted,
		InternalEncrypted:       user.InternalEncrypted,
		Metadata:                user.Metadata,
		CreatedAt:               ts.New(user.CreatedAt),
		UpdatedAt:               ts.New(user.UpdatedAt),
		EncryptedAt:             ts.New(user.EncryptedAt),
		ExternalEncryptionLevel: int32(user.ExternalEncryptionLevel),
		InternalEncryptionLevel: int32(user.InternalEncryptionLevel),
	}
}

func ProtoUserToUser(user *go_block.User) *User {
	if user == nil {
		return nil
	}
	return &User{
		Id:                      user.Id,
		OptionalId:              user.OptionalId,
		Email:                   user.Email,
		Password:                user.Password,
		Image:                   user.Image,
		ExternalEncrypted:       user.ExternalEncrypted,
		InternalEncrypted:       user.InternalEncrypted,
		Metadata:                user.Metadata,
		CreatedAt:               user.CreatedAt.AsTime(),
		UpdatedAt:               user.UpdatedAt.AsTime(),
		EncryptedAt:             user.EncryptedAt.AsTime(),
		ExternalEncryptionLevel: int(user.ExternalEncryptionLevel),
		InternalEncryptionLevel: int(user.InternalEncryptionLevel),
	}
}
