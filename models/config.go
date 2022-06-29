package models

import (
	"github.com/nuntiodev/hera-sdks/go_hera"
	"github.com/nuntiodev/x/cryptox"
	ts "google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type Config struct {
	Id                       string              `bson:"_id" json:"id"`
	Name                     cryptox.Stringx     `bson:"name" json:"name"`
	Logo                     cryptox.Stringx     `bson:"logo" json:"logo"`
	NuntioVerifyId           cryptox.Stringx     `bson:"nuntio_verify_id" json:"nuntio_verify_id"`
	DisableSignup            bool                `bson:"disable_signup" json:"disable_signup"`
	DisableLogin             bool                `bson:"disable_login" json:"disable_login"`
	CreatedAt                time.Time           `bson:"created_at" json:"created_at"`
	UpdatedAt                time.Time           `bson:"updated_at" json:"updated_at"`
	ValidatePassword         bool                `bson:"validate_password" json:"validate_password"`
	VerifyEmail              bool                `bson:"verify_email" json:"verify_email"`
	VerifyPhone              bool                `bson:"verify_phone" json:"verify_phone"`
	SupportedLoginMechanisms []go_hera.LoginType `bson:"login_type" json:"login_type"`
	PublicKey                string              `bson:"public_key" json:"public_key"`
}

func ProtoConfigToConfig(config *go_hera.Config) *Config {
	if config == nil {
		return nil
	}
	return &Config{
		Name:                     cryptox.Stringx{Body: config.Name},
		Logo:                     cryptox.Stringx{Body: config.Logo},
		DisableSignup:            config.DisableSignup,
		DisableLogin:             config.DisableLogin,
		CreatedAt:                config.CreatedAt.AsTime(),
		UpdatedAt:                config.UpdatedAt.AsTime(),
		ValidatePassword:         config.ValidatePassword,
		SupportedLoginMechanisms: config.SupportedLoginMechanisms,
		VerifyEmail:              config.VerifyEmail,
		VerifyPhone:              config.VerifyPhone,
		PublicKey:                config.PublicKey,
	}
}

func ConfigToProtoConfig(config *Config) *go_hera.Config {
	if config == nil {
		return nil
	}
	return &go_hera.Config{
		Name:                     config.Name.Body,
		Logo:                     config.Logo.Body,
		DisableSignup:            config.DisableSignup,
		DisableLogin:             config.DisableLogin,
		CreatedAt:                ts.New(config.CreatedAt),
		UpdatedAt:                ts.New(config.UpdatedAt),
		ValidatePassword:         config.ValidatePassword,
		SupportedLoginMechanisms: config.SupportedLoginMechanisms,
		VerifyEmail:              config.VerifyEmail,
		VerifyPhone:              config.VerifyPhone,
		PublicKey:                config.PublicKey,
	}
}
