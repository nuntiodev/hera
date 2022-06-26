package models

import (
	"fmt"
	"github.com/nuntiodev/hera-proto/go_hera"
)

const (
	loginTypeEmailPassword    = "email/password"
	loginTypeEmailCode        = "email/code"
	loginTypePhonePassword    = "phone/password"
	loginTypePhoneCode        = "phone/code"
	loginTypeUsernamePassword = "username/password"
)

type HeraConfig struct {
	AppName          string   `bson:"app_name" json:"app_name"`
	Logo             string   `bson:"logo" json:"logo"`
	DisableSignup    bool     `bson:"disable_signup" json:"disable_signup"`
	DisableLogin     bool     `bson:"disable_login" json:"disable_login"`
	ValidatePassword bool     `bson:"validate_password" json:"validate_password"`
	VerifyEmail      bool     `bson:"verify_email" json:"verify_email"`
	VerifyPhone      bool     `bson:"verify_phone" json:"verify_phone"`
	LoginMechanisms  []string `bson:"login_mechanisms" json:"login_mechanisms"`
	PublicKey        string   `bson:"public_key" json:"public_key"`
}

func getSupportedLoginMechanisms(m []string) ([]go_hera.LoginType, error) {
	var resp []go_hera.LoginType
	for _, loginType := range m {
		switch loginType {
		case loginTypeEmailPassword:
			resp = append(resp, go_hera.LoginType_EMAIL_PASSWORD)
		case loginTypeEmailCode:
			resp = append(resp, go_hera.LoginType_EMAIL_VERIFICATION_CODE)
		case loginTypePhonePassword:
			resp = append(resp, go_hera.LoginType_PHONE_PASSWORD)
		case loginTypePhoneCode:
			resp = append(resp, go_hera.LoginType_PHONE_VERIFICATION_CODE)
		case loginTypeUsernamePassword:
			resp = append(resp, go_hera.LoginType_USERNAME_PASSWORD)
		default:
			return nil, fmt.Errorf("invalid login type: %s \n Supported login types are: %s, %s, %s, %s, %s, %s", m, loginTypeEmailPassword, loginTypeEmailCode, loginTypePhonePassword, loginTypePhoneCode, loginTypeUsernamePassword)
		}
	}
	return resp, nil
}

func HeraConfigToProto(h *HeraConfig) (*go_hera.Config, error) {
	if h == nil {
		return nil, nil
	}
	supportedLoginMechanisms, err := getSupportedLoginMechanisms(h.LoginMechanisms)
	if err != nil {
		return nil, err
	}
	return &go_hera.Config{
		Name:                     h.AppName,
		Logo:                     h.Logo,
		DisableSignup:            h.DisableSignup,
		DisableLogin:             h.DisableLogin,
		ValidatePassword:         h.ValidatePassword,
		VerifyPhone:              h.VerifyPhone,
		VerifyEmail:              h.VerifyEmail,
		PublicKey:                h.PublicKey,
		SupportedLoginMechanisms: supportedLoginMechanisms,
	}, nil
}
