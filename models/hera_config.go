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

type HeraApp struct {
	Name             string   `bson:"name" json:"name"`
	Logo             string   `bson:"logo" json:"logo"`
	DisableSignup    bool     `bson:"disable_signup" json:"disable_signup"`
	DisableLogin     bool     `bson:"disable_login" json:"disable_login"`
	ValidatePassword bool     `bson:"validate_password" json:"validate_password"`
	VerifyEmail      bool     `bson:"verify_email" json:"verify_email"`
	VerifyPhone      bool     `bson:"verify_phone" json:"verify_phone"`
	LoginMechanisms  []string `bson:"login_mechanisms" json:"login_mechanisms"`
	PublicKey        string   `bson:"public_key" json:"public_key"`
}

type HeraUser struct {
	FirstName string `bson:"first_name" json:"first_name"`
	LastName  string `bson:"last_name" json:"last_name"`
	Email     string `bson:"email" json:"email"`
	Password  string `bson:"password" json:"password"`
	Phone     string `bson:"phone" json:"phone"`
	Username  string `bson:"username" json:"username"`
	Image     string `bson:"image" json:"image"`
	Id        string `bson:"id" json:"id"`
}

type HeraConfig struct {
	App   HeraApp    `bson:"app" json:"app"`
	Users []HeraUser `bson:"users" json:"users"`
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

func HeraConfigToProtoConfig(h *HeraConfig) (*go_hera.Config, error) {
	if h == nil {
		return nil, nil
	}
	supportedLoginMechanisms, err := getSupportedLoginMechanisms(h.App.LoginMechanisms)
	if err != nil {
		return nil, err
	}
	return &go_hera.Config{
		Name:                     h.App.Name,
		Logo:                     h.App.Logo,
		DisableSignup:            h.App.DisableSignup,
		DisableLogin:             h.App.DisableLogin,
		ValidatePassword:         h.App.ValidatePassword,
		VerifyPhone:              h.App.VerifyPhone,
		VerifyEmail:              h.App.VerifyEmail,
		PublicKey:                h.App.PublicKey,
		SupportedLoginMechanisms: supportedLoginMechanisms,
	}, nil
}

func HeraConfigToProtoUsers(h *HeraConfig) []*go_hera.User {
	if h == nil {
		return nil
	}
	var resp []*go_hera.User
	for _, user := range h.Users {
		resp = append(resp, &go_hera.User{
			FirstName: &user.FirstName,
			LastName:  &user.LastName,
			Email:     &user.Email,
			Password:  user.Password,
			Phone:     &user.Phone,
			Username:  &user.Username,
			Image:     &user.Image,
			Id:        user.Id,
		})
	}
	return resp
}
