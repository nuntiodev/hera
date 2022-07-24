package models

import (
	"github.com/nuntiodev/hera-sdks/go_hera"
	"github.com/nuntiodev/x/cryptox"
	ts "google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type Config struct {
	Id                       string                  `bson:"_id" json:"id"`
	Name                     cryptox.Stringx         `bson:"name" json:"name"`
	Logo                     cryptox.Stringx         `bson:"logo" json:"logo"`
	NuntioVerifyId           cryptox.Stringx         `bson:"nuntio_verify_id" json:"nuntio_verify_id"`
	DisableSignup            bool                    `bson:"disable_signup" json:"disable_signup"`
	DisableLogin             bool                    `bson:"disable_login" json:"disable_login"`
	CreatedAt                time.Time               `bson:"created_at" json:"created_at"`
	UpdatedAt                time.Time               `bson:"updated_at" json:"updated_at"`
	ValidatePassword         bool                    `bson:"validate_password" json:"validate_password"`
	VerifyEmail              bool                    `bson:"verify_email" json:"verify_email"`
	VerifyPhone              bool                    `bson:"verify_phone" json:"verify_phone"`
	SupportedLoginMechanisms []go_hera.LoginType     `bson:"login_type" json:"login_type"`
	PublicKey                cryptox.Stringx         `bson:"public_key" json:"public_key"`
	HashingAlgorithm         go_hera.HasingAlgorithm `bson:"hashing_algorithm" json:"hashing_algorithm"`
	Bcrypt                   *Bcrypt                 `bson:"bcrypt" json:"bcrypt"`
	Scrypt                   *Scrypt                 `bson:"scrypt" json:"scrypt"`
}

func ProtoConfigToConfig(config *go_hera.Config) *Config {
	if config == nil {
		return nil
	}
	var bcrypt *Bcrypt
	if config.Bcrypt != nil {
		bcrypt = &Bcrypt{
			Cost: int(config.Bcrypt.Cost),
		}
	}
	var scrypt *Scrypt
	if config.Scrypt != nil {
		scrypt = &Scrypt{
			SignerKey:     cryptox.Stringx{Body: config.Scrypt.SignerKey},
			SaltSeparator: cryptox.Stringx{Body: config.Scrypt.SaltSeparator},
			Rounds:        int(config.Scrypt.Rounds),
			MemCost:       int(config.Scrypt.MemCost),
			P:             int(config.Scrypt.P),
			KeyLen:        int(config.Scrypt.KeyLen),
		}
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
		PublicKey:                cryptox.Stringx{Body: config.PublicKey},
		HashingAlgorithm:         config.HasingAlgorithm,
		Bcrypt:                   bcrypt,
		Scrypt:                   scrypt,
	}
}

func ConfigToProtoConfig(config *Config) *go_hera.Config {
	if config == nil {
		return nil
	}
	var bcrypt *go_hera.Bcrypt
	if config.Bcrypt != nil {
		bcrypt = &go_hera.Bcrypt{
			Cost: int32(config.Bcrypt.Cost),
		}
	}
	var scrypt *go_hera.Scrypt
	if config.Scrypt != nil {
		scrypt = &go_hera.Scrypt{
			SignerKey:     config.Scrypt.SignerKey.Body,
			SaltSeparator: config.Scrypt.SaltSeparator.Body,
			Rounds:        int32(config.Scrypt.Rounds),
			MemCost:       int32(config.Scrypt.MemCost),
			P:             int32(config.Scrypt.P),
			KeyLen:        int32(config.Scrypt.KeyLen),
		}
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
		PublicKey:                config.PublicKey.Body,
		HasingAlgorithm:          config.HashingAlgorithm,
		Bcrypt:                   bcrypt,
		Scrypt:                   scrypt,
	}
}
