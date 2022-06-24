package models

import (
	"github.com/araddon/dateparse"
	"github.com/nuntiodev/hera-proto/go_hera"
	"github.com/nuntiodev/x/cryptox"
	ts "google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type User struct {
	Id                         string               `bson:"_id" json:"id"`
	Username                   cryptox.Stringx      `bson:"username" json:"username"`
	Email                      cryptox.Stringx      `bson:"email" json:"email"`
	Password                   string               `bson:"password" json:"password"`
	Image                      cryptox.Stringx      `bson:"image" json:"image"`
	Metadata                   cryptox.Stringx      `bson:"metadata" json:"metadata"`
	CreatedAt                  time.Time            `bson:"created_at" json:"created_at"`
	UpdatedAt                  time.Time            `bson:"updated_at" json:"updated_at"`
	FirstName                  cryptox.Stringx      `bson:"first_name" json:"first_name"`
	LastName                   cryptox.Stringx      `bson:"last_name" json:"last_name"`
	Birthdate                  cryptox.Stringx      `bson:"birthdate" json:"birthdate"`
	VerificationEmailSentAt    time.Time            `bson:"verification_email_sent_at" json:"verification_email_sent_at"`
	EmailVerificationCode      string               `bson:"email_verification_code" json:"email_verification_code"`
	VerificationEmailExpiresAt time.Time            `bson:"verification_email_expires_at" json:"verification_email_expires_at"`
	VerifyEmailAttempts        int32                `bson:"verify_email_attempts" json:"verify_email_attempts"`
	ResetPasswordCode          string               `bson:"reset_password_code" json:"reset_password_code"`
	ResetPasswordCodeSentAt    time.Time            `bson:"reset_password_code_sent_at" json:"reset_password_code_sent_at"`
	ResetPasswordCodeExpiresAt time.Time            `bson:"reset_password_code_expires_at" json:"reset_password_code_expires_at"`
	ResetPasswordAttempts      int32                `bson:"reset_password_attempts" json:"reset_password_attempts"`
	VerifiedEmails             []string             `bson:"verified_emails" json:"verified_emails"`
	EmailHash                  string               `bson:"email_hash" json:"email_hash"`
	Phone                      cryptox.Stringx      `bson:"phone" json:"phone"`
	PhoneHash                  string               `bson:"phone_hash" json:"phone_hash"`
	VerificationTextSentAt     time.Time            `bson:"verification_text_sent_at" json:"verification_text_sent_at"`
	VerifiedPhoneNumbers       []string             `bson:"verified_phone_numbers" json:"verified_phone_numbers"`
	PreferredLanguage          go_hera.LanguageCode `bson:"preferred_language" json:"preferred_language"`
	UsernameHash               string               `bson:"username_hash" json:"username_hash"`
	VerifyPhoneAttempts        int32                `bson:"verify_phone_attempts" json:"verify_phone_attempts"`
	PhoneVerificationCode      string               `bson:"phone_verification_code" json:"phone_verification_code"`
}

func UsersToProto(users []*User) []*go_hera.User {
	var resp []*go_hera.User
	for _, user := range users {
		resp = append(resp, UserToProtoUser(user))
	}
	return resp
}

func UserToProtoUser(user *User) *go_hera.User {
	if user == nil {
		return nil
	}
	birthdate := &ts.Timestamp{}
	if user.Birthdate.Body != "" {
		t, err := dateparse.ParseAny(user.Birthdate.Body)
		if err == nil {
			birthdate = ts.New(t)
		}
	}
	return &go_hera.User{
		Id:                         user.Id,
		Username:                   &user.Username.Body,
		Email:                      &user.Email.Body,
		Password:                   user.Password,
		Image:                      &user.Image.Body,
		FirstName:                  &user.FirstName.Body,
		LastName:                   &user.LastName.Body,
		Birthdate:                  birthdate,
		Metadata:                   user.Metadata.Body,
		CreatedAt:                  ts.New(user.CreatedAt),
		UpdatedAt:                  ts.New(user.UpdatedAt),
		VerificationEmailSentAt:    ts.New(user.VerificationEmailSentAt),
		EmailVerificationCode:      user.EmailVerificationCode,
		VerificationEmailExpiresAt: ts.New(user.VerificationEmailExpiresAt),
		ResetPasswordCode:          user.ResetPasswordCode,
		ResetPasswordCodeSentAt:    ts.New(user.ResetPasswordCodeSentAt),
		ResetPasswordCodeExpiresAt: ts.New(user.ResetPasswordCodeExpiresAt),
		ResetPasswordAttempts:      user.ResetPasswordAttempts,
		VerifyEmailAttempts:        user.VerifyEmailAttempts,
		VerifiedEmails:             user.VerifiedEmails,
		EmailHash:                  user.EmailHash,
		Phone:                      &user.Phone.Body,
		PhoneHash:                  user.PhoneHash,
		VerificationTextSentAt:     ts.New(user.VerificationTextSentAt),
		VerifiedPhoneNumbers:       user.VerifiedPhoneNumbers,
		PreferredLanguage:          &user.PreferredLanguage,
		UsernameHash:               user.UsernameHash,
		VerifyPhoneAttempts:        user.VerifyPhoneAttempts,
		PhoneVerificationCode:      user.PhoneVerificationCode,
	}
}

func ProtoUserToUser(user *go_hera.User) *User {
	if user == nil {
		return nil
	}
	birthdate := ""
	if user.Birthdate != nil {
		birthdate = user.GetBirthdate().AsTime().String()
	}
	return &User{
		Id:                         user.Id,
		Username:                   cryptox.Stringx{Body: user.GetUsername()},
		Email:                      cryptox.Stringx{Body: user.GetEmail()},
		Password:                   user.Password,
		Image:                      cryptox.Stringx{Body: user.GetImage()},
		FirstName:                  cryptox.Stringx{Body: user.GetFirstName()},
		LastName:                   cryptox.Stringx{Body: user.GetLastName()},
		Birthdate:                  cryptox.Stringx{Body: birthdate},
		Metadata:                   cryptox.Stringx{Body: user.GetMetadata()},
		CreatedAt:                  user.GetCreatedAt().AsTime(),
		UpdatedAt:                  user.GetUpdatedAt().AsTime(),
		VerificationEmailSentAt:    user.GetVerificationEmailSentAt().AsTime(),
		EmailVerificationCode:      user.GetEmailVerificationCode(),
		VerificationEmailExpiresAt: user.GetVerificationEmailExpiresAt().AsTime(),
		ResetPasswordCode:          user.GetResetPasswordCode(),
		ResetPasswordCodeSentAt:    user.GetResetPasswordCodeSentAt().AsTime(),
		ResetPasswordCodeExpiresAt: user.GetResetPasswordCodeExpiresAt().AsTime(),
		ResetPasswordAttempts:      user.GetResetPasswordAttempts(),
		VerifyEmailAttempts:        user.GetVerifyEmailAttempts(),
		VerifiedEmails:             user.GetVerifiedEmails(),
		EmailHash:                  user.GetEmailHash(),
		Phone:                      cryptox.Stringx{Body: user.GetPhone()},
		PhoneHash:                  user.GetPhoneHash(),
		VerificationTextSentAt:     user.GetVerificationTextSentAt().AsTime(),
		VerifiedPhoneNumbers:       user.GetVerifiedPhoneNumbers(),
		PreferredLanguage:          user.GetPreferredLanguage(),
		UsernameHash:               user.GetUsernameHash(),
		VerifyPhoneAttempts:        user.VerifyPhoneAttempts,
		PhoneVerificationCode:      user.PhoneVerificationCode,
	}
}
