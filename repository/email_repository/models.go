package email_repository

import (
	"github.com/nuntiodev/block-proto/go_block"
	ts "google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type Email struct {
	Id                      string    `bson:"_id" json:"_id"`
	Logo                    string    `bson:"logo" json:"logo"`
	WelcomeMessage          string    `bson:"welcome_message" json:"welcome_message"`
	BodyMessage             string    `bson:"body_message" json:"body_message"`
	FooterMessage           string    `bson:"footer_message" json:"footer_message"`
	CreatedAt               time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt               time.Time `bson:"updated_at" json:"updated_at"`
	EncryptedAt             time.Time `bson:"encrypted_at" json:"encrypted_at"`
	TriggerOnCreate         bool      `bson:"trigger_on_create" json:"trigger_on_create"`
	InternalEncryptionLevel int32     `bson:"internal_encryption_level" json:"internal_encryption_level"`
	Subject                 string    `bson:"subject" json:"subject"`
	TemplatePath            string    `bson:"template_path" json:"template_path"`
}

func EmailToProtoEmail(email *Email) *go_block.Email {
	if email == nil {
		return nil
	}
	return &go_block.Email{
		Id:                      email.Id,
		Logo:                    email.Logo,
		WelcomeMessage:          email.WelcomeMessage,
		BodyMessage:             email.BodyMessage,
		FooterMessage:           email.FooterMessage,
		CreatedAt:               ts.New(email.CreatedAt),
		UpdatedAt:               ts.New(email.UpdatedAt),
		EncryptedAt:             ts.New(email.EncryptedAt),
		TriggerOnCreate:         email.TriggerOnCreate,
		InternalEncryptionLevel: email.InternalEncryptionLevel,
		Subject:                 email.Subject,
		TemplatePath:            email.TemplatePath,
	}
}

func ProtoEmailToEmail(email *go_block.Email) *Email {
	if email == nil {
		return nil
	}
	return &Email{
		Id:                      email.Id,
		Logo:                    email.Logo,
		WelcomeMessage:          email.WelcomeMessage,
		BodyMessage:             email.BodyMessage,
		FooterMessage:           email.FooterMessage,
		CreatedAt:               email.CreatedAt.AsTime(),
		UpdatedAt:               email.UpdatedAt.AsTime(),
		EncryptedAt:             email.EncryptedAt.AsTime(),
		TriggerOnCreate:         email.TriggerOnCreate,
		InternalEncryptionLevel: email.InternalEncryptionLevel,
		Subject:                 email.Subject,
		TemplatePath:            email.TemplatePath,
	}
}
