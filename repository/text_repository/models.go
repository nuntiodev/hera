package text_repository

import "github.com/nuntiodev/block-proto/go_block"

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

type ProfileText struct {
	ProfileTitle              string `bson:"profile_title" json:"profile_title"`
	Logout                    string `bson:"logout" json:"logout"`
	ChangeEmailTitle          string `bson:"change_email_title" json:"change_email_title"`
	ChangeEmailDescription    string `bson:"change_email_description" json:"change_email_description"`
	ChangePasswordTitle       string `bson:"change_password_title" json:"change_password_title"`
	ChangePasswordDescription string `bson:"change_password_description" json:"change_password_description"`
}

type Text struct {
	Id                      go_block.LanguageCode `bson:"_id" json:"id"`
	GeneralText             *GeneralText          `bson:"general_text" json:"general_text"`
	WelcomeText             *WelcomeText          `bson:"welcome_text" json:"welcome_text"`
	LoginText               *LoginText            `bson:"login_text" json:"login_text"`
	RegisterText            *RegisterText         `bson:"register_text" json:"register_text"`
	ProfileText             *ProfileText          `bson:"profile_text" json:"profile_text"`
	InternalEncryptionLevel int32                 `bson:"internal_encryption_level" json:"internal_encryption_level"`
}

func ProtoGeneralTextToGeneralText(text *go_block.GeneralText) *GeneralText {
	if text == nil {
		return nil
	}
	return &GeneralText{
		MissingPasswordTitle:   text.MissingPasswordTitle,
		MissingPasswordDetails: text.MissingPasswordDetails,
		MissingEmailTitle:      text.MissingEmailTitle,
		MissingEmailDetails:    text.MissingEmailDetails,
		CreatedBy:              text.CreatedBy,
		PasswordHint:           text.PasswordHint,
		EmailHint:              text.EmailHint,
		ErrorTitle:             text.ErrorTitle,
		ErrorDescription:       text.ErrorDescription,
		NoWifiTitle:            text.NoWifiTitle,
		NoWifiDescription:      text.NoWifiDescription,
	}
}

func GeneralTextToProtoGeneralText(text *GeneralText) *go_block.GeneralText {
	if text == nil {
		return nil
	}
	return &go_block.GeneralText{
		MissingPasswordTitle:   text.MissingPasswordTitle,
		MissingPasswordDetails: text.MissingPasswordDetails,
		MissingEmailTitle:      text.MissingEmailTitle,
		MissingEmailDetails:    text.MissingEmailDetails,
		CreatedBy:              text.CreatedBy,
		PasswordHint:           text.PasswordHint,
		EmailHint:              text.EmailHint,
		ErrorTitle:             text.ErrorTitle,
		ErrorDescription:       text.ErrorDescription,
		NoWifiTitle:            text.NoWifiTitle,
		NoWifiDescription:      text.NoWifiDescription,
	}
}

func ProtoWelcomeTextToWelcomeText(text *go_block.WelcomeText) *WelcomeText {
	if text == nil {
		return nil
	}
	return &WelcomeText{
		WelcomeTitle:       text.WelcomeTitle,
		WelcomeDetails:     text.WelcomeDetails,
		ContinueWithNuntio: text.ContinueWithNuntio,
	}
}

func WelcomeTextToProtoWelcomeText(text *WelcomeText) *go_block.WelcomeText {
	if text == nil {
		return nil
	}
	return &go_block.WelcomeText{
		WelcomeTitle:       text.WelcomeTitle,
		WelcomeDetails:     text.WelcomeDetails,
		ContinueWithNuntio: text.ContinueWithNuntio,
	}
}

func ProtoLoginTextToLoginText(text *go_block.LoginText) *LoginText {
	if text == nil {
		return nil
	}
	return &LoginText{
		LoginButton:    text.LoginButton,
		LoginTitle:     text.LoginTitle,
		LoginDetails:   text.LoginDetails,
		ForgotPassword: text.ForgotPassword,
	}
}

func LoginTextToProtoLoginText(text *LoginText) *go_block.LoginText {
	if text == nil {
		return nil
	}
	return &go_block.LoginText{
		LoginButton:    text.LoginButton,
		LoginTitle:     text.LoginTitle,
		LoginDetails:   text.LoginDetails,
		ForgotPassword: text.ForgotPassword,
	}
}

func ProtoRegisterTextToRegisterText(text *go_block.RegisterText) *RegisterText {
	if text == nil {
		return nil
	}
	return &RegisterText{
		RegisterButton:            text.RegisterButton,
		RegisterTitle:             text.RegisterTitle,
		RegisterDetails:           text.RegisterDetails,
		PasswordDoNotMatchTitle:   text.PasswordDoNotMatchTitle,
		PasswordDoNotMatchDetails: text.PasswordDoNotMatchDetails,
		RepeatPasswordHint:        text.RepeatPasswordHint,
		ContainsSpecialChar:       text.ContainsSpecialChar,
		ContainsNumberChar:        text.ContainsNumberChar,
		PasswordMustMatch:         text.PasswordMustMatch,
		ContainsEightChars:        text.ContainsEightChars,
	}
}

func RegisterTextToProtoRegisterText(text *RegisterText) *go_block.RegisterText {
	if text == nil {
		return nil
	}
	return &go_block.RegisterText{
		RegisterButton:            text.RegisterButton,
		RegisterTitle:             text.RegisterTitle,
		RegisterDetails:           text.RegisterDetails,
		PasswordDoNotMatchTitle:   text.PasswordDoNotMatchTitle,
		PasswordDoNotMatchDetails: text.PasswordDoNotMatchDetails,
		RepeatPasswordHint:        text.RepeatPasswordHint,
		ContainsSpecialChar:       text.ContainsSpecialChar,
		ContainsNumberChar:        text.ContainsNumberChar,
		PasswordMustMatch:         text.PasswordMustMatch,
		ContainsEightChars:        text.ContainsEightChars,
	}
}

func ProtoProfileTextToProfileText(text *go_block.ProfileText) *ProfileText {
	if text == nil {
		return nil
	}
	return &ProfileText{
		ProfileTitle:              text.ProfileTitle,
		Logout:                    text.Logout,
		ChangeEmailTitle:          text.ChangeEmailTitle,
		ChangeEmailDescription:    text.ChangeEmailDescription,
		ChangePasswordTitle:       text.ChangePasswordTitle,
		ChangePasswordDescription: text.ChangePasswordDescription,
	}
}

func ProfileTextToProtoProfileText(text *ProfileText) *go_block.ProfileText {
	if text == nil {
		return nil
	}
	return &go_block.ProfileText{
		ProfileTitle:              text.ProfileTitle,
		Logout:                    text.Logout,
		ChangeEmailTitle:          text.ChangeEmailTitle,
		ChangeEmailDescription:    text.ChangeEmailDescription,
		ChangePasswordTitle:       text.ChangePasswordTitle,
		ChangePasswordDescription: text.ChangePasswordDescription,
	}
}

func ProtoTextToText(text *go_block.Text) *Text {
	if text == nil {
		return nil
	}
	return &Text{
		Id:           text.Id,
		GeneralText:  ProtoGeneralTextToGeneralText(text.GeneralText),
		WelcomeText:  ProtoWelcomeTextToWelcomeText(text.WelcomeText),
		LoginText:    ProtoLoginTextToLoginText(text.LoginText),
		RegisterText: ProtoRegisterTextToRegisterText(text.RegisterText),
		ProfileText:  ProtoProfileTextToProfileText(text.ProfileText),
	}
}

func TextToProtoText(text *Text) *go_block.Text {
	if text == nil {
		return nil
	}
	return &go_block.Text{
		Id:                      text.Id,
		GeneralText:             GeneralTextToProtoGeneralText(text.GeneralText),
		WelcomeText:             WelcomeTextToProtoWelcomeText(text.WelcomeText),
		LoginText:               LoginTextToProtoLoginText(text.LoginText),
		RegisterText:            RegisterTextToProtoRegisterText(text.RegisterText),
		ProfileText:             ProfileTextToProtoProfileText(text.ProfileText),
		InternalEncryptionLevel: text.InternalEncryptionLevel,
	}
}
