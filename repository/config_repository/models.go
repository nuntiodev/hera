package config_repository

import (
	"github.com/nuntiodev/block-proto/go_block"
	ts "google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type AuthConfig struct {
	Logo                      string `bson:"logo" json:"logo"`
	WelcomeTitle              string `bson:"welcome_title" json:"welcome_title"`
	WelcomeDetails            string `bson:"welcome_details" json:"welcome_details"`
	LoginButton               string `bson:"login_button" json:"login_button"`
	LoginTitle                string `bson:"login_title" json:"login_title"`
	LoginDetails              string `bson:"login_details" json:"login_details"`
	RegisterButton            string `bson:"register_button" json:"register_button"`
	RegisterTitle             string `bson:"register_title" json:"register_title"`
	RegisterDetails           string `bson:"register_details" json:"register_details"`
	MissingPasswordTitle      string `bson:"missing_password_title" json:"missing_password_title"`
	MissingPasswordDetails    string `bson:"missing_password_details" json:"missing_password_details"`
	MissingEmailTitle         string `bson:"missing_email_title" json:"missing_email_title"`
	MissingEmailDetails       string `bson:"missing_email_details" json:"missing_email_details"`
	PasswordDoNotMatchTitle   string `bson:"password_do_not_match_title" json:"password_do_not_match_title"`
	PasswordDoNotMatchDetails string `bson:"password_do_not_match_details" json:"password_do_not_match_details"`
	CreatedBy                 string `bson:"created_by" json:"created_by"`
}

type Config struct {
	Id                      string      `bson:"_id" json:"id"`
	Name                    string      `bson:"name" json:"name"`
	Website                 string      `bson:"website" json:"website"`
	About                   string      `bson:"about" json:"about"`
	Email                   string      `bson:"email" json:"email"`
	Logo                    string      `bson:"logo" json:"logo"`
	Terms                   string      `bson:"terms" json:"terms"`
	EnableNuntioConnect     bool        `bson:"enable_nuntio_connect" json:"enable_nuntio_connect"`
	DisableDefaultSignup    bool        `bson:"disable_default_signup" json:"disable_default_signup"`
	DisableDefaultLogin     bool        `bson:"disable_default_login" json:"disable_default_login"`
	AuthConfig              *AuthConfig `bson:"auth_config" json:"auth_config"`
	CreatedAt               time.Time   `bson:"created_at" json:"created_at"`
	UpdatedAt               time.Time   `bson:"updated_at" json:"updated_at"`
	InternalEncryptionLevel int32       `bson:"internal_encryption_level" json:"internal_encryption_level"`
}

func ProtoConfigToConfig(config *go_block.Config) *Config {
	if config == nil {
		return nil
	}
	authConfig := &AuthConfig{}
	if config.AuthConfig != nil {
		authConfig = &AuthConfig{
			Logo:                      config.AuthConfig.Logo,
			WelcomeTitle:              config.AuthConfig.WelcomeTitle,
			WelcomeDetails:            config.AuthConfig.WelcomeDetails,
			LoginButton:               config.AuthConfig.LoginButton,
			LoginTitle:                config.AuthConfig.LoginTitle,
			LoginDetails:              config.AuthConfig.LoginDetails,
			RegisterButton:            config.AuthConfig.RegisterButton,
			RegisterTitle:             config.AuthConfig.RegisterTitle,
			RegisterDetails:           config.AuthConfig.RegisterDetails,
			MissingPasswordTitle:      config.AuthConfig.MissingPasswordTitle,
			MissingPasswordDetails:    config.AuthConfig.MissingPasswordDetails,
			MissingEmailTitle:         config.AuthConfig.MissingEmailTitle,
			MissingEmailDetails:       config.AuthConfig.MissingEmailDetails,
			PasswordDoNotMatchTitle:   config.AuthConfig.PasswordDoNotMatchTitle,
			PasswordDoNotMatchDetails: config.AuthConfig.PasswordDoNotMatchDetails,
			CreatedBy:                 config.AuthConfig.CreatedBy,
		}
	}
	return &Config{
		Id:                      config.Id,
		Name:                    config.Name,
		Website:                 config.Website,
		About:                   config.About,
		Email:                   config.Email,
		Logo:                    config.Logo,
		Terms:                   config.Terms,
		EnableNuntioConnect:     config.EnableNuntioConnect,
		DisableDefaultSignup:    config.DisableDefaultSignup,
		DisableDefaultLogin:     config.DisableDefaultLogin,
		AuthConfig:              authConfig,
		CreatedAt:               config.CreatedAt.AsTime(),
		UpdatedAt:               config.UpdatedAt.AsTime(),
		InternalEncryptionLevel: config.InternalEncryptionLevel,
	}
}

func ConfigToProtoConfig(config *Config) *go_block.Config {
	if config == nil {
		return nil
	}
	authConfig := &go_block.AuthConfig{}
	if config.AuthConfig != nil {
		authConfig = &go_block.AuthConfig{
			Logo:                      config.AuthConfig.Logo,
			WelcomeTitle:              config.AuthConfig.WelcomeTitle,
			WelcomeDetails:            config.AuthConfig.WelcomeDetails,
			LoginButton:               config.AuthConfig.LoginButton,
			LoginTitle:                config.AuthConfig.LoginTitle,
			LoginDetails:              config.AuthConfig.LoginDetails,
			RegisterButton:            config.AuthConfig.RegisterButton,
			RegisterTitle:             config.AuthConfig.RegisterTitle,
			RegisterDetails:           config.AuthConfig.RegisterDetails,
			MissingPasswordTitle:      config.AuthConfig.MissingPasswordTitle,
			MissingPasswordDetails:    config.AuthConfig.MissingPasswordDetails,
			MissingEmailTitle:         config.AuthConfig.MissingEmailTitle,
			MissingEmailDetails:       config.AuthConfig.MissingEmailDetails,
			PasswordDoNotMatchTitle:   config.AuthConfig.PasswordDoNotMatchTitle,
			PasswordDoNotMatchDetails: config.AuthConfig.PasswordDoNotMatchDetails,
			CreatedBy:                 config.AuthConfig.CreatedBy,
		}
	}
	return &go_block.Config{
		Id:                      config.Id,
		Name:                    config.Name,
		Website:                 config.Website,
		About:                   config.About,
		Email:                   config.Email,
		Logo:                    config.Logo,
		Terms:                   config.Terms,
		EnableNuntioConnect:     config.EnableNuntioConnect,
		DisableDefaultSignup:    config.DisableDefaultSignup,
		DisableDefaultLogin:     config.DisableDefaultLogin,
		AuthConfig:              authConfig,
		CreatedAt:               ts.New(config.CreatedAt),
		UpdatedAt:               ts.New(config.UpdatedAt),
		InternalEncryptionLevel: config.InternalEncryptionLevel,
	}
}
