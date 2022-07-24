package models

import (
	"fmt"
	"github.com/nuntiodev/hera-sdks/go_hera"
	"github.com/nuntiodev/x/pointerx"
	ts "google.golang.org/protobuf/types/known/timestamppb"
	"strings"
	"time"
)

const (
	loginTypeEmailPassword    = "email/password"
	loginTypeEmailCode        = "email/code"
	loginTypePhonePassword    = "phone/password"
	loginTypePhoneCode        = "phone/code"
	loginTypeUsernamePassword = "username/password"
)

type HeraApp struct {
	Name             string      `bson:"name" json:"name"`
	Logo             string      `bson:"logo" json:"logo"`
	DisableSignup    bool        `bson:"disable_signup" json:"disable_signup"`
	DisableLogin     bool        `bson:"disable_login" json:"disable_login"`
	ValidatePassword bool        `bson:"validate_password" json:"validate_password"`
	VerifyEmail      bool        `bson:"verify_email" json:"verify_email"`
	VerifyPhone      bool        `bson:"verify_phone" json:"verify_phone"`
	LoginMechanisms  []string    `bson:"login_mechanisms" json:"login_mechanisms"`
	PublicKey        string      `bson:"public_key" json:"public_key"`
	HashingAlgorithm string      `bson:"hashing_algorithm" json:"hashing_algorithm"`
	Scrypt           *ScryptHera `bson:"scrypt" json:"scrypt"`
	Bcrypt           *Bcrypt     `bson:"bcrypt" json:"bcrypt"`
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
	Birthdate string `bson:"birthdate" json:"birthdate"`
}

type HeraConfig struct {
	App   HeraApp    `bson:"default_config" json:"default_config"`
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
	hashingAlgorithm := go_hera.HasingAlgorithm_BCRYPT
	if strings.ToLower(strings.TrimSpace(h.App.HashingAlgorithm)) == "scrypt" {
		hashingAlgorithm = go_hera.HasingAlgorithm_SCRYPT
	}
	var bcrypt *go_hera.Bcrypt
	if h.App.Bcrypt != nil {
		bcrypt = &go_hera.Bcrypt{Cost: int32(h.App.Bcrypt.Cost)}
	}
	var scrypt *go_hera.Scrypt
	if h.App.Scrypt != nil {
		scrypt = &go_hera.Scrypt{
			SignerKey:     h.App.Scrypt.SignerKey,
			SaltSeparator: h.App.Scrypt.SaltSeparator,
			Rounds:        int32(h.App.Scrypt.Rounds),
			MemCost:       int32(h.App.Scrypt.MemCost),
			P:             int32(h.App.Scrypt.P),
			KeyLen:        int32(h.App.Scrypt.KeyLen),
		}
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
		HasingAlgorithm:          hashingAlgorithm,
		Bcrypt:                   bcrypt,
		Scrypt:                   scrypt,
		SupportedLoginMechanisms: supportedLoginMechanisms,
	}, nil
}

func HeraConfigToProtoUsers(h *HeraConfig) []*go_hera.User {
	if h == nil {
		return nil
	}
	var resp []*go_hera.User
	for _, user := range h.Users {
		birthdate, _ := time.Parse("2006-01-02", user.Birthdate)
		resp = append(resp, &go_hera.User{
			FirstName: pointerx.StringPtr(user.FirstName),
			LastName:  pointerx.StringPtr(user.LastName),
			Email:     pointerx.StringPtr(user.Email),
			Password:  &go_hera.Hash{Body: user.Password},
			Phone:     pointerx.StringPtr(user.Phone),
			Username:  pointerx.StringPtr(user.Username),
			Image:     pointerx.StringPtr(user.Image),
			Id:        user.Id,
			Birthdate: ts.New(birthdate),
		})
	}
	return resp
}
