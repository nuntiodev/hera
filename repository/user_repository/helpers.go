package user_repository

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/nuntiodev/x/pointerx"
	"go.mongodb.org/mongo-driver/bson"
	"net/mail"
	"regexp"
	"strings"
	"unicode"

	"github.com/google/uuid"
	"github.com/nuntiodev/hera-sdks/go_hera"
	ts "google.golang.org/protobuf/types/known/timestamppb"
)

var (
	ErrInvalidUpdatedAt = errors.New("invalid updated at timestamp")
)

func prepare(action int, user *go_hera.User) {
	if user == nil {
		return
	}
	switch action {
	case actionCreate:
		user.CreatedAt = ts.Now()
		user.UpdatedAt = ts.Now()
		user.VerificationEmailSentAt = &ts.Timestamp{}
		user.VerifiedEmails = []string{}
		if strings.TrimSpace(user.Id) == "" {
			user.Id = uuid.NewString()
		}
	}
	user.Id = strings.TrimSpace(user.Id)
	if user.Email != nil {
		user.Email = pointerx.StringPtr(strings.TrimSpace(strings.ToLower(*user.Email)))
	}
	if user.FirstName != nil {
		user.FirstName = pointerx.StringPtr(strings.TrimSpace(*user.FirstName))
	}
	if user.LastName != nil {
		user.LastName = pointerx.StringPtr(strings.TrimSpace(*user.LastName))
	}
	if user.Image != nil {
		user.Image = pointerx.StringPtr(strings.TrimSpace(*user.Image))
	}
	if user.Username != nil {
		user.Username = pointerx.StringPtr(strings.TrimSpace(*user.Username))
	}
	if user.Phone != nil {
		user.Phone = pointerx.StringPtr(strings.TrimSpace(*user.Phone))
	}
	user.Metadata = strings.TrimSpace(user.Metadata)
	user.EmailVerificationCode = strings.TrimSpace(user.EmailVerificationCode)
	user.PhoneVerificationCode = strings.TrimSpace(user.PhoneVerificationCode)
	user.EmailHash = strings.TrimSpace(user.EmailHash)
	user.PhoneHash = strings.TrimSpace(user.PhoneHash)
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

func validateEmail(email string) error {
	_, err := mail.ParseAddress(email)
	if err != nil && email != "" {
		return err
	}
	return nil
}

func validatePhone(phone string) error {
	re := regexp.MustCompile(`^(?:(?:\(?(?:00|\+)([1-4]\d\d|[1-9]\d?)\)?)?[\-\.\ \\\/]?)?((?:\(?\d{1,}\)?[\-\.\ \\\/]?){0,})(?:[\-\.\ \\\/]?(?:#|ext\.?|extension|x)[\-\.\ \\\/]?(\d+))?$`)
	if re.MatchString(phone) == false && phone != "" {
		return errors.New("invalid phone number")
	}
	return nil
}

func validateMetadata(metadata string) error {
	if !json.Valid([]byte(metadata)) && metadata != "" {
		return errors.New("invalid json data")
	}
	return nil
}

func getUserFilter(user *go_hera.User) (bson.M, error) {
	if user == nil {
		return nil, UserIsNilErr
	}
	if user.Id != "" {
		return bson.M{"_id": user.Id}, nil
	}
	emailHash, usernameHash, phoneHash := generateUserHashes(user)
	if emailHash != "" {
		return bson.M{"email_hash": emailHash}, nil
	}
	if usernameHash != "" {
		return bson.M{"username_hash": usernameHash}, nil
	}
	if phoneHash != "" {
		return bson.M{"phone_number_hash": phoneHash}, nil
	}
	return nil, errors.New("no identifier")
}

func generateUserHashes(user *go_hera.User) (emailHash string, usernameHash string, phoneHash string) {
	if user == nil {
		return
	}
	if user.GetEmail() != "" {
		emailHash = fmt.Sprintf("%x", md5.Sum([]byte(user.GetEmail())))
	}
	if user.GetUsername() != "" {
		usernameHash = fmt.Sprintf("%x", md5.Sum([]byte(user.GetUsername())))
	}
	if user.GetPhone() != "" {
		phoneHash = fmt.Sprintf("%x", md5.Sum([]byte(user.GetPhone())))
	}
	return
}
