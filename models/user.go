package models

import (
	"github.com/araddon/dateparse"
	"github.com/nuntiodev/block-proto/go_block"
	"github.com/nuntiodev/x/cryptox"
	ts "google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type User struct {
	Id                          string                `bson:"_id" json:"id"`
	Username                    cryptox.Stringx       `bson:"username" json:"username"`
	Email                       cryptox.Stringx       `bson:"email" json:"email"`
	Password                    string                `bson:"password" json:"password"`
	Image                       cryptox.Stringx       `bson:"image" json:"image"`
	Metadata                    cryptox.Stringx       `bson:"metadata" json:"metadata"`
	CreatedAt                   time.Time             `bson:"created_at" json:"created_at"`
	UpdatedAt                   time.Time             `bson:"updated_at" json:"updated_at"`
	EncryptedAt                 time.Time             `bson:"encrypted_at" json:"encrypted_at"`
	FirstName                   cryptox.Stringx       `bson:"first_name" json:"first_name"`
	LastName                    cryptox.Stringx       `bson:"last_name" json:"last_name"`
	Birthdate                   cryptox.Stringx       `bson:"birthdate" json:"birthdate"`
	EmailVerifiedAt             time.Time             `bson:"email_verified_at" json:"email_verified_at"`
	EmailIsVerified             bool                  `bson:"email_is_verified" json:"email_is_verified"`
	VerificationEmailSentAt     time.Time             `bson:"verification_email_sent_at" json:"verification_email_sent_at"`
	EmailVerificationCode       string                `bson:"email_verification_code" json:"email_verification_code"`
	VerificationEmailExpiresAt  time.Time             `bson:"verification_email_expires_at" json:"verification_email_expires_at"`
	VerifyEmailAttempts         int32                 `bson:"verify_email_attempts" json:"verify_email_attempts"`
	ResetPasswordCode           string                `bson:"reset_password_code" json:"reset_password_code"`
	ResetPasswordEmailSentAt    time.Time             `bson:"reset_password_email_sent_at" json:"reset_password_email_sent_at"`
	ResetPasswordEmailExpiresAt time.Time             `bson:"reset_password_email_expires_at" json:"reset_password_email_expires_at"`
	ResetPasswordAttempts       int32                 `bson:"reset_password_attempts" json:"reset_password_attempts"`
	VerifiedEmails              []string              `bson:"verified_emails" json:"verified_emails"`
	EmailHash                   string                `bson:"email_hash" json:"email_hash"`
	PhoneNumber                 cryptox.Stringx       `bson:"phone_number" json:"phone_number"`
	PhoneNumberHash             string                `bson:"phone_number_hash" json:"phone_number_hash"`
	PhoneNumberIsVerified       bool                  `bson:"phone_number_is_verified" json:"phone_number_is_verified"`
	VerificationTextSentAt      time.Time             `bson:"verification_text_sent_at" json:"verification_text_sent_at"`
	VerifiedPhoneNumbers        []string              `bson:"verified_phone_numbers" json:"verified_phone_numbers"`
	PreferredLanguage           go_block.LanguageCode `bson:"preferred_language" json:"preferred_language"`
	UsernameHash                string                `bson:"username_hash" json:"username_hash"`
}

func UsersToProto(users []*User) []*go_block.User {
	var resp []*go_block.User
	for _, user := range users {
		resp = append(resp, UserToProtoUser(user))
	}
	return resp
}

func UserToProtoUser(user *User) *go_block.User {
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
	return &go_block.User{
		Id:                          user.Id,
		Username:                    user.Username.Body,
		Email:                       user.Email.Body,
		Password:                    user.Password,
		Image:                       user.Image.Body,
		FirstName:                   user.FirstName.Body,
		LastName:                    user.LastName.Body,
		Birthdate:                   birthdate,
		Metadata:                    user.Metadata.Body,
		CreatedAt:                   ts.New(user.CreatedAt),
		UpdatedAt:                   ts.New(user.UpdatedAt),
		EncryptedAt:                 ts.New(user.EncryptedAt),
		VerificationEmailSentAt:     ts.New(user.VerificationEmailSentAt),
		EmailVerifiedAt:             ts.New(user.EmailVerifiedAt),
		EmailIsVerified:             user.EmailIsVerified,
		EmailVerificationCode:       user.EmailVerificationCode,
		VerificationEmailExpiresAt:  ts.New(user.VerificationEmailExpiresAt),
		ResetPasswordCode:           user.ResetPasswordCode,
		ResetPasswordEmailSentAt:    ts.New(user.ResetPasswordEmailSentAt),
		ResetPasswordEmailExpiresAt: ts.New(user.ResetPasswordEmailExpiresAt),
		ResetPasswordAttempts:       user.ResetPasswordAttempts,
		VerifyEmailAttempts:         user.VerifyEmailAttempts,
		VerifiedEmails:              user.VerifiedEmails,
		EmailHash:                   user.EmailHash,
		PhoneNumber:                 user.PhoneNumber.Body,
		PhoneNumberHash:             user.PhoneNumberHash,
		PhoneNumberIsVerified:       user.PhoneNumberIsVerified,
		VerificationTextSentAt:      ts.New(user.VerificationTextSentAt),
		VerifiedPhoneNumbers:        user.VerifiedPhoneNumbers,
		PreferredLanguage:           user.PreferredLanguage,
		UsernameHash:                user.UsernameHash,
	}
}

func ProtoUserToUser(user *go_block.User) *User {
	if user == nil {
		return nil
	}
	birthdate := ""
	if user.Birthdate != nil {
		birthdate = user.Birthdate.AsTime().String()
	}
	return &User{
		Id:                          user.Id,
		Username:                    cryptox.Stringx{Body: user.Username},
		Email:                       cryptox.Stringx{Body: user.Email},
		Password:                    user.Password,
		Image:                       cryptox.Stringx{Body: user.Image},
		FirstName:                   cryptox.Stringx{Body: user.FirstName},
		LastName:                    cryptox.Stringx{Body: user.LastName},
		Birthdate:                   cryptox.Stringx{Body: birthdate},
		Metadata:                    cryptox.Stringx{Body: user.Metadata},
		CreatedAt:                   user.CreatedAt.AsTime(),
		UpdatedAt:                   user.UpdatedAt.AsTime(),
		EncryptedAt:                 user.EncryptedAt.AsTime(),
		VerificationEmailSentAt:     user.VerificationEmailSentAt.AsTime(),
		EmailVerifiedAt:             user.EmailVerifiedAt.AsTime(),
		EmailIsVerified:             user.EmailIsVerified,
		EmailVerificationCode:       user.EmailVerificationCode,
		VerificationEmailExpiresAt:  user.VerificationEmailExpiresAt.AsTime(),
		ResetPasswordCode:           user.ResetPasswordCode,
		ResetPasswordEmailSentAt:    user.ResetPasswordEmailSentAt.AsTime(),
		ResetPasswordEmailExpiresAt: user.ResetPasswordEmailExpiresAt.AsTime(),
		ResetPasswordAttempts:       user.ResetPasswordAttempts,
		VerifyEmailAttempts:         user.VerifyEmailAttempts,
		VerifiedEmails:              user.VerifiedEmails,
		EmailHash:                   user.EmailHash,
		PhoneNumber:                 cryptox.Stringx{Body: user.PhoneNumber},
		PhoneNumberHash:             user.PhoneNumberHash,
		PhoneNumberIsVerified:       user.PhoneNumberIsVerified,
		VerificationTextSentAt:      user.VerificationTextSentAt.AsTime(),
		VerifiedPhoneNumbers:        user.VerifiedPhoneNumbers,
		PreferredLanguage:           user.PreferredLanguage,
		UsernameHash:                user.UsernameHash,
	}
}
