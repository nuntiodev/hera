package user_repository

import (
	"github.com/nuntiodev/block-proto/go_block"
	ts "google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type User struct {
	Id                          string    `bson:"_id" json:"id"`
	OptionalId                  string    `bson:"optional_id" json:"optional_id"`
	Email                       string    `bson:"email" json:"email"`
	EmailHash                   string    `bson:"email_hash" json:"email_hash"`
	Password                    string    `bson:"password" json:"password"`
	Image                       string    `bson:"image" json:"image"`
	InternalEncryptionLevel     int       `bson:"internal_encryption_level" json:"internal_encryption_level"`
	ExternalEncryptionLevel     int       `bson:"external_encryption_level" json:"external_encryption_level"`
	Metadata                    string    `bson:"metadata" json:"metadata"`
	CreatedAt                   time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt                   time.Time `bson:"updated_at" json:"updated_at"`
	EncryptedAt                 time.Time `bson:"encrypted_at" json:"encrypted_at"`
	RequireEmailVerification    bool      `bson:"require_email_verification" json:"require_email_verification"`
	FirstName                   string    `bson:"first_name" json:"first_name"`
	LastName                    string    `bson:"last_name" json:"last_name"`
	Birthdate                   string    `bson:"birthdate" json:"birthdate"`
	EmailVerifiedAt             time.Time `bson:"email_verified_at" json:"email_verified_at"`
	EmailIsVerified             bool      `bson:"email_is_verified" json:"email_is_verified"`
	VerificationEmailSentAt     time.Time `bson:"verification_email_sent_at" json:"verification_email_sent_at"`
	EmailVerificationCode       string    `bson:"email_verification_code" json:"email_verification_code"`
	VerificationEmailExpiresAt  time.Time `bson:"verification_email_expires_at" json:"verification_email_expires_at"`
	VerifyEmailAttempts         int32     `bson:"verify_email_attempts" json:"verify_email_attempts"`
	ResetPasswordCode           string    `bson:"reset_password_code" json:"reset_password_code"`
	ResetPasswordEmailSentAt    time.Time `bson:"reset_password_email_sent_at" json:"reset_password_email_sent_at"`
	ResetPasswordEmailExpiresAt time.Time `bson:"reset_password_email_expires_at" json:"reset_password_email_expires_at"`
	ResetPasswordAttempts       int32     `bson:"reset_password_attempts" json:"reset_password_attempts"`
	VerifiedEmails              []string  `bson:"verified_emails" json:"verified_emails"`
}

func UserToProtoUser(user *User) *go_block.User {
	if user == nil {
		return nil
	}
	birthdate := &ts.Timestamp{}
	if user.Birthdate != "" {
		t, err := time.Parse(time.RFC3339, user.Birthdate)
		if err == nil {
			birthdate = ts.New(t)
		}
	}
	return &go_block.User{
		Id:                          user.Id,
		OptionalId:                  user.OptionalId,
		Email:                       user.Email,
		Password:                    user.Password,
		Image:                       user.Image,
		FirstName:                   user.FirstName,
		LastName:                    user.LastName,
		Birthdate:                   birthdate,
		Metadata:                    user.Metadata,
		CreatedAt:                   ts.New(user.CreatedAt),
		UpdatedAt:                   ts.New(user.UpdatedAt),
		EncryptedAt:                 ts.New(user.EncryptedAt),
		ExternalEncryptionLevel:     int32(user.ExternalEncryptionLevel),
		InternalEncryptionLevel:     int32(user.InternalEncryptionLevel),
		VerificationEmailSentAt:     ts.New(user.VerificationEmailSentAt),
		EmailVerifiedAt:             ts.New(user.EmailVerifiedAt),
		RequireEmailVerification:    user.RequireEmailVerification,
		EmailIsVerified:             user.EmailIsVerified,
		EmailVerificationCode:       user.EmailVerificationCode,
		VerificationEmailExpiresAt:  ts.New(user.VerificationEmailExpiresAt),
		ResetPasswordCode:           user.ResetPasswordCode,
		ResetPasswordEmailSentAt:    ts.New(user.ResetPasswordEmailSentAt),
		ResetPasswordEmailExpiresAt: ts.New(user.ResetPasswordEmailExpiresAt),
		ResetPasswordAttempts:       user.ResetPasswordAttempts,
		VerifyEmailAttempts:         user.VerifyEmailAttempts,
		VerifiedEmails:              user.VerifiedEmails,
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
		OptionalId:                  user.OptionalId,
		Email:                       user.Email,
		Password:                    user.Password,
		Image:                       user.Image,
		FirstName:                   user.FirstName,
		LastName:                    user.LastName,
		Birthdate:                   birthdate,
		Metadata:                    user.Metadata,
		CreatedAt:                   user.CreatedAt.AsTime(),
		UpdatedAt:                   user.UpdatedAt.AsTime(),
		EncryptedAt:                 user.EncryptedAt.AsTime(),
		ExternalEncryptionLevel:     int(user.ExternalEncryptionLevel),
		InternalEncryptionLevel:     int(user.InternalEncryptionLevel),
		VerificationEmailSentAt:     user.VerificationEmailSentAt.AsTime(),
		EmailVerifiedAt:             user.EmailVerifiedAt.AsTime(),
		EmailIsVerified:             user.EmailIsVerified,
		RequireEmailVerification:    user.RequireEmailVerification,
		EmailVerificationCode:       user.EmailVerificationCode,
		VerificationEmailExpiresAt:  user.VerificationEmailExpiresAt.AsTime(),
		ResetPasswordCode:           user.ResetPasswordCode,
		ResetPasswordEmailSentAt:    user.ResetPasswordEmailSentAt.AsTime(),
		ResetPasswordEmailExpiresAt: user.ResetPasswordEmailExpiresAt.AsTime(),
		ResetPasswordAttempts:       user.ResetPasswordAttempts,
		VerifyEmailAttempts:         user.VerifyEmailAttempts,
		VerifiedEmails:              user.VerifiedEmails,
	}
}
