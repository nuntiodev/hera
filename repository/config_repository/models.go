package config_repository

import (
	"github.com/nuntiodev/block-proto/go_block"
	ts "google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type GeneralText struct {
	MissingPasswordTitle   string `bson:"missing_password_title" json:"missing_password_title"`
	MissingPasswordDetails string `bson:"missing_password_details" json:"missing_password_details"`
	MissingEmailTitle      string `bson:"missing_email_title" json:"missing_email_title"`
	MissingEmailDetails    string `bson:"missing_email_details" json:"missing_email_details"`
	CreatedBy              string `bson:"created_by" json:"created_by"`
	PasswordHint           string `bson:"password_hint" json:"password_hint"`
	EmailHint              string `bson:"email_hint" json:"email_hint"`
	ErrorTitle             string `bson:"error_title" json:"error_title"`
	ErrorDescription       string `bson:"error_description" json:"error_description"`
	NoWifiTitle            string `bson:"no_wifi_title" json:"no_wifi_title"`
	NoWifiDescription      string `bson:"no_wifi_description" json:"no_wifi_description"`
}

type WelcomeText struct {
	WelcomeTitle       string `bson:"welcome_title" json:"welcome_title"`
	WelcomeDetails     string `bson:"welcome_details" json:"welcome_details"`
	ContinueWithNuntio string `bson:"continue_with_nuntio" json:"continue_with_nuntio"`
}

type RegisterText struct {
	RegisterButton            string `bson:"register_button" json:"register_button"`
	RegisterTitle             string `bson:"register_title" json:"register_title"`
	RegisterDetails           string `bson:"register_details" json:"register_details"`
	PasswordDoNotMatchTitle   string `bson:"password_do_not_match_title" json:"password_do_not_match_title"`
	PasswordDoNotMatchDetails string `bson:"password_do_not_match_details" json:"password_do_not_match_details"`
	RepeatPasswordHint        string `bson:"repeat_password_hint" json:"repeat_password_hint"`
	ContainsSpecialChar       string `bson:"contains_special_char" json:"contains_special_char"`
	ContainsNumberChar        string `bson:"contains_number_char" json:"contains_number_char"`
	PasswordMustMatch         string `bson:"password_must_match" json:"password_must_match"`
	ContainsEightChars        string `bson:"contains_eight_chars" json:"contains_eight_chars"`
}

type LoginText struct {
	LoginButton    string `bson:"login_button" json:"login_button"`
	LoginTitle     string `bson:"login_title" json:"login_title"`
	LoginDetails   string `bson:"login_details" json:"login_details"`
	ForgotPassword string `bson:"forgot_password" json:"forgot_password"`
}

type Config struct {
	Id                       string        `bson:"_id" json:"id"`
	Name                     string        `bson:"name" json:"name"`
	Logo                     string        `bson:"logo" json:"logo"`
	EnableNuntioConnect      bool          `bson:"enable_nuntio_connect" json:"enable_nuntio_connect"`
	DisableDefaultSignup     bool          `bson:"disable_default_signup" json:"disable_default_signup"`
	DisableDefaultLogin      bool          `bson:"disable_default_login" json:"disable_default_login"`
	GeneralText              *GeneralText  `bson:"general_text" json:"general_text"`
	WelcomeText              *WelcomeText  `bson:"welcome_text" json:"welcome_text"`
	LoginText                *LoginText    `bson:"login_text" json:"login_text"`
	RegisterText             *RegisterText `bson:"register_text" json:"register_text"`
	CreatedAt                time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt                time.Time     `bson:"updated_at" json:"updated_at"`
	InternalEncryptionLevel  int32         `bson:"internal_encryption_level" json:"internal_encryption_level"`
	ValidatePassword         bool          `bson:"validate_password" json:"validate_password"`
	NuntioConnectId          string        `bson:"nuntio_connect_id" json:"nuntio_connect_id"`
	RequireEmailVerification bool          `bson:"require_email_verification" json:"require_email_verification"`
}

func ProtoConfigToConfig(config *go_block.Config) *Config {
	if config == nil {
		return nil
	}
	generalText := &GeneralText{}
	if config.GeneralText != nil {
		generalText = &GeneralText{
			MissingPasswordTitle:   config.GeneralText.MissingPasswordTitle,
			MissingPasswordDetails: config.GeneralText.MissingPasswordDetails,
			MissingEmailTitle:      config.GeneralText.MissingEmailTitle,
			MissingEmailDetails:    config.GeneralText.MissingEmailDetails,
			CreatedBy:              config.GeneralText.CreatedBy,
			PasswordHint:           config.GeneralText.PasswordHint,
			EmailHint:              config.GeneralText.EmailHint,
			ErrorTitle:             config.GeneralText.ErrorTitle,
			ErrorDescription:       config.GeneralText.ErrorDescription,
			NoWifiTitle:            config.GeneralText.NoWifiTitle,
			NoWifiDescription:      config.GeneralText.NoWifiDescription,
		}
	}
	welcomeText := &WelcomeText{}
	if config.WelcomeText != nil {
		welcomeText = &WelcomeText{
			WelcomeTitle:   config.WelcomeText.WelcomeTitle,
			WelcomeDetails: config.WelcomeText.WelcomeDetails,
		}
	}
	registerText := &RegisterText{}
	if config.RegisterText != nil {
		registerText = &RegisterText{
			RegisterButton:            config.RegisterText.RegisterButton,
			RegisterTitle:             config.RegisterText.RegisterTitle,
			RegisterDetails:           config.RegisterText.RegisterDetails,
			PasswordDoNotMatchTitle:   config.RegisterText.PasswordDoNotMatchTitle,
			PasswordDoNotMatchDetails: config.RegisterText.PasswordDoNotMatchDetails,
			RepeatPasswordHint:        config.RegisterText.RepeatPasswordHint,
			ContainsSpecialChar:       config.RegisterText.ContainsSpecialChar,
			ContainsNumberChar:        config.RegisterText.ContainsNumberChar,
			PasswordMustMatch:         config.RegisterText.PasswordMustMatch,
			ContainsEightChars:        config.RegisterText.ContainsEightChars,
		}
	}
	loginText := &LoginText{}
	if config.LoginText != nil {
		loginText = &LoginText{
			LoginButton:    config.LoginText.LoginButton,
			LoginTitle:     config.LoginText.LoginTitle,
			LoginDetails:   config.LoginText.LoginDetails,
			ForgotPassword: config.LoginText.ForgotPassword,
		}
	}
	return &Config{
		Id:                       config.Id,
		Name:                     config.Name,
		Logo:                     config.Logo,
		EnableNuntioConnect:      config.EnableNuntioConnect,
		DisableDefaultSignup:     config.DisableDefaultSignup,
		DisableDefaultLogin:      config.DisableDefaultLogin,
		GeneralText:              generalText,
		WelcomeText:              welcomeText,
		RegisterText:             registerText,
		LoginText:                loginText,
		CreatedAt:                config.CreatedAt.AsTime(),
		UpdatedAt:                config.UpdatedAt.AsTime(),
		InternalEncryptionLevel:  config.InternalEncryptionLevel,
		ValidatePassword:         config.ValidatePassword,
		NuntioConnectId:          config.NuntioConnectId,
		RequireEmailVerification: config.RequireEmailVerification,
	}
}

func ConfigToProtoConfig(config *Config) *go_block.Config {
	if config == nil {
		return nil
	}
	generalText := &go_block.GeneralText{}
	if config.GeneralText != nil {
		generalText = &go_block.GeneralText{
			MissingPasswordTitle:   config.GeneralText.MissingPasswordTitle,
			MissingPasswordDetails: config.GeneralText.MissingPasswordDetails,
			MissingEmailTitle:      config.GeneralText.MissingEmailTitle,
			MissingEmailDetails:    config.GeneralText.MissingEmailDetails,
			CreatedBy:              config.GeneralText.CreatedBy,
			PasswordHint:           config.GeneralText.PasswordHint,
			EmailHint:              config.GeneralText.EmailHint,
			ErrorTitle:             config.GeneralText.ErrorTitle,
			ErrorDescription:       config.GeneralText.ErrorDescription,
			NoWifiTitle:            config.GeneralText.NoWifiTitle,
			NoWifiDescription:      config.GeneralText.NoWifiDescription,
		}
	}
	welcomeText := &go_block.WelcomeText{}
	if config.GeneralText != nil {
		welcomeText = &go_block.WelcomeText{
			WelcomeTitle:   config.WelcomeText.WelcomeTitle,
			WelcomeDetails: config.WelcomeText.WelcomeDetails,
		}
	}
	registerText := &go_block.RegisterText{}
	if config.RegisterText != nil {
		registerText = &go_block.RegisterText{
			RegisterButton:            config.RegisterText.RegisterButton,
			RegisterTitle:             config.RegisterText.RegisterTitle,
			RegisterDetails:           config.RegisterText.RegisterDetails,
			PasswordDoNotMatchTitle:   config.RegisterText.PasswordDoNotMatchTitle,
			PasswordDoNotMatchDetails: config.RegisterText.PasswordDoNotMatchDetails,
			RepeatPasswordHint:        config.RegisterText.RepeatPasswordHint,
			ContainsSpecialChar:       config.RegisterText.ContainsSpecialChar,
			ContainsNumberChar:        config.RegisterText.ContainsNumberChar,
			PasswordMustMatch:         config.RegisterText.PasswordMustMatch,
			ContainsEightChars:        config.RegisterText.ContainsEightChars,
		}
	}
	loginText := &go_block.LoginText{}
	if config.LoginText != nil {
		loginText = &go_block.LoginText{
			LoginButton:    config.LoginText.LoginButton,
			LoginTitle:     config.LoginText.LoginTitle,
			LoginDetails:   config.LoginText.LoginDetails,
			ForgotPassword: config.LoginText.ForgotPassword,
		}
	}
	return &go_block.Config{
		Id:                       config.Id,
		Name:                     config.Name,
		Logo:                     config.Logo,
		EnableNuntioConnect:      config.EnableNuntioConnect,
		DisableDefaultSignup:     config.DisableDefaultSignup,
		DisableDefaultLogin:      config.DisableDefaultLogin,
		GeneralText:              generalText,
		WelcomeText:              welcomeText,
		RegisterText:             registerText,
		LoginText:                loginText,
		CreatedAt:                ts.New(config.CreatedAt),
		UpdatedAt:                ts.New(config.UpdatedAt),
		InternalEncryptionLevel:  config.InternalEncryptionLevel,
		ValidatePassword:         config.ValidatePassword,
		NuntioConnectId:          config.NuntioConnectId,
		RequireEmailVerification: config.RequireEmailVerification,
	}
}
