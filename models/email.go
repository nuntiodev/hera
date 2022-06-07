package models

import (
	"github.com/nuntiodev/block-proto/go_block"
	"github.com/nuntiodev/x/cryptox"
	ts "google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type Email struct {
	Id              string                `bson:"_id" json:"_id"`
	Logo            cryptox.Stringx       `bson:"logo" json:"logo"`
	WelcomeMessage  cryptox.Stringx       `bson:"welcome_message" json:"welcome_message"`
	BodyMessage     cryptox.Stringx       `bson:"body_message" json:"body_message"`
	FooterMessage   cryptox.Stringx       `bson:"footer_message" json:"footer_message"`
	CreatedAt       time.Time             `bson:"created_at" json:"created_at"`
	UpdatedAt       time.Time             `bson:"updated_at" json:"updated_at"`
	EncryptedAt     time.Time             `bson:"encrypted_at" json:"encrypted_at"`
	TriggerOnCreate bool                  `bson:"trigger_on_create" json:"trigger_on_create"`
	Subject         cryptox.Stringx       `bson:"subject" json:"subject"`
	TemplatePath    cryptox.Stringx       `bson:"template_path" json:"template_path"`
	LanguageCode    go_block.LanguageCode `bson:"language_code" json:"language_code"`
}

func EmailToProtoEmail(email *Email) *go_block.Email {
	if email == nil {
		return nil
	}
	return &go_block.Email{
		Id:              email.Id,
		Logo:            email.Logo.Body,
		WelcomeMessage:  email.WelcomeMessage.Body,
		BodyMessage:     email.BodyMessage.Body,
		FooterMessage:   email.FooterMessage.Body,
		CreatedAt:       ts.New(email.CreatedAt),
		UpdatedAt:       ts.New(email.UpdatedAt),
		TriggerOnCreate: email.TriggerOnCreate,
		Subject:         email.Subject.Body,
		TemplatePath:    email.TemplatePath.Body,
		LanguageCode:    email.LanguageCode,
	}
}

func ProtoEmailToEmail(email *go_block.Email) *Email {
	if email == nil {
		return nil
	}
	return &Email{
		Id:              email.Id,
		Logo:            cryptox.Stringx{Body: email.Logo},
		WelcomeMessage:  cryptox.Stringx{Body: email.WelcomeMessage},
		BodyMessage:     cryptox.Stringx{Body: email.BodyMessage},
		FooterMessage:   cryptox.Stringx{Body: email.FooterMessage},
		CreatedAt:       email.CreatedAt.AsTime(),
		UpdatedAt:       email.UpdatedAt.AsTime(),
		TriggerOnCreate: email.TriggerOnCreate,
		Subject:         cryptox.Stringx{Body: email.Subject},
		TemplatePath:    cryptox.Stringx{Body: email.TemplatePath},
		LanguageCode:    email.LanguageCode,
	}
}
