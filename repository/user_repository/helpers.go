package user_repository

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"strings"
	"unicode"

	"github.com/badoux/checkmail"
	"github.com/nuntiodev/block-proto/go_block"
	uuid "github.com/satori/go.uuid"
	ts "google.golang.org/protobuf/types/known/timestamppb"
)

func prepare(action int, user *go_block.User) {
	if user == nil {
		return
	}
	switch action {
	case actionCreate:
		user.CreatedAt = ts.Now()
		user.UpdatedAt = ts.Now()
		user.VerificationEmailSentAt = &ts.Timestamp{}
		user.EmailVerifiedAt = &ts.Timestamp{}
		user.EmailIsVerified = false
		if strings.TrimSpace(user.Id) == "" {
			user.Id = uuid.NewV4().String()
		}
	case actionUpdatePassword, actionUpdateImage, actionUpdateMetadata,
		actionUpdateNamespace, actionUpdateSecurity, actionUpdateEmail,
		actionUpdateOptionalId, actionUpdateName, actionUpdateBirthdate,
		actionUpdateEmailVerified, actionUpdateVerificationEmailSent:
		user.UpdatedAt = ts.Now()
	}
	user.Id = strings.TrimSpace(user.Id)
	user.Email = strings.TrimSpace(strings.ToLower(user.Email))
	user.FirstName = strings.TrimSpace(user.FirstName)
	user.LastName = strings.TrimSpace(user.LastName)
	user.Image = strings.TrimSpace(user.Image)
	user.OptionalId = strings.TrimSpace(user.OptionalId)
	user.Metadata = strings.TrimSpace(user.Metadata)
	user.VerificationCode = strings.TrimSpace(user.VerificationCode)
}

func (r *mongodbRepository) validate(action int, user *go_block.User) error {
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
	case actionUpdateSecurity, actionUpdateEmailVerified:
		if !user.UpdatedAt.IsValid() {
			return errors.New("invalid updated at")
		}
	case actionUpdateVerificationEmailSent:
		if user.Id == "" && user.Email == "" && user.OptionalId == "" {
			return errors.New("missing required search parameter")
		} else if user.VerificationCode == "" {
			return errors.New("missing required verification code")
		}
	case actionGetAll, actionUpdateOptionalId:
		return nil
	}
	if len(user.Email) > maxFieldLength {
		return errors.New("email field is too long")
	} else if len(user.OptionalId) > maxFieldLength {
		return errors.New("optional id field is too long")
	} else if len(user.Metadata) > 10*maxFieldLength {
		return errors.New("metadata field is too long")
	} else if len(user.FirstName) > maxFieldLength {
		return errors.New("first name field is too long")
	} else if len(user.LastName) > maxFieldLength {
		return errors.New("last name field is too long")
	}
	return nil
}

func validatePassword(password string) error {
	var (
		num, sym bool
		tot      uint8
	)
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			tot++
		case unicode.IsLower(char):
			tot++
		case unicode.IsNumber(char):
			num = true
			tot++
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			sym = true
			tot++
		default:
			continue
		}
	}
	if !num {
		return errors.New("missing required number")
	} else if !sym {
		return errors.New("missing required symbol")
	} else if tot < 8 {
		return errors.New("password is too short; must be at least 8 chars long")
	}
	return nil
}

func getUserFilter(user *go_block.User) (bson.M, error) {
	if user == nil {
		return nil, errors.New("user is nil")
	}
	filter := bson.M{}
	if user.Id != "" {
		filter = bson.M{"_id": user.Id}
	} else if user.Email != "" {
		filter = bson.M{"email_hash": fmt.Sprintf("%x", md5.Sum([]byte(user.Email)))}
		fmt.Println(fmt.Sprintf("%x", md5.Sum([]byte(user.Email))))
	} else if user.OptionalId != "" {
		filter = bson.M{"optional_id": user.OptionalId}
	} else {
		return nil, errors.New("missing search filter")
	}
	return filter, nil
}
