package models

import (
	"github.com/nuntiodev/block-proto/go_block"
	"github.com/nuntiodev/x/cryptox"
	ts "google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type Config struct {
	Id                             string                `bson:"_id" json:"id"`
	Name                           cryptox.Stringx       `bson:"name" json:"name"`
	Logo                           cryptox.Stringx       `bson:"logo" json:"logo"`
	EnableNuntioConnect            bool                  `bson:"enable_nuntio_connect" json:"enable_nuntio_connect"`
	DisableDefaultSignup           bool                  `bson:"disable_default_signup" json:"disable_default_signup"`
	DisableDefaultLogin            bool                  `bson:"disable_default_login" json:"disable_default_login"`
	CreatedAt                      time.Time             `bson:"created_at" json:"created_at"`
	UpdatedAt                      time.Time             `bson:"updated_at" json:"updated_at"`
	ValidatePassword               bool                  `bson:"validate_password" json:"validate_password"`
	NuntioConnectId                string                `bson:"nuntio_connect_id" json:"nuntio_connect_id"`
	RequireEmailVerification       bool                  `bson:"require_email_verification" json:"require_email_verification"`
	LoginType                      go_block.LoginType    `bson:"login_type" json:"login_type"`
	RequirePhoneNumberVerification bool                  `bson:"require_phone_number_verification" json:"require_phone_number_verification"`
	DefaultLanguage                go_block.LanguageCode `bson:"default_language" json:"default_language"`
}

func ProtoConfigToConfig(config *go_block.Config) *Config {
	if config == nil {
		return nil
	}
	return &Config{
		Id:                             config.Id,
		Name:                           cryptox.Stringx{Body: config.Name},
		Logo:                           cryptox.Stringx{Body: config.Logo},
		EnableNuntioConnect:            config.EnableNuntioConnect,
		DisableDefaultSignup:           config.DisableDefaultSignup,
		DisableDefaultLogin:            config.DisableDefaultLogin,
		CreatedAt:                      config.CreatedAt.AsTime(),
		UpdatedAt:                      config.UpdatedAt.AsTime(),
		ValidatePassword:               config.ValidatePassword,
		NuntioConnectId:                config.NuntioConnectId,
		RequireEmailVerification:       config.RequireEmailVerification,
		LoginType:                      config.LoginType,
		RequirePhoneNumberVerification: config.RequirePhoneNumberVerification,
	}
}

func ConfigToProtoConfig(config *Config) *go_block.Config {
	if config == nil {
		return nil
	}
	return &go_block.Config{
		Id:                             config.Id,
		Name:                           config.Name.Body,
		Logo:                           config.Logo.Body,
		EnableNuntioConnect:            config.EnableNuntioConnect,
		DisableDefaultSignup:           config.DisableDefaultSignup,
		DisableDefaultLogin:            config.DisableDefaultLogin,
		CreatedAt:                      ts.New(config.CreatedAt),
		UpdatedAt:                      ts.New(config.UpdatedAt),
		ValidatePassword:               config.ValidatePassword,
		NuntioConnectId:                config.NuntioConnectId,
		RequireEmailVerification:       config.RequireEmailVerification,
		LoginType:                      config.LoginType,
		RequirePhoneNumberVerification: config.RequirePhoneNumberVerification,
	}
}