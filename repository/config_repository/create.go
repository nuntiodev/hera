package config_repository

import (
	"context"
	"errors"
	"github.com/nuntiodev/block-proto/go_block"
)

func (cr *defaultConfigRepository) Create(ctx context.Context, config *go_block.Config) (*go_block.Config, error) {
	prepare(actionCreate, config)
	if config == nil {
		return nil, errors.New("missing required config")
	} else if config.Id == "" {
		return nil, errors.New("missing required id")
	}
	create := ProtoConfigToConfig(config)
	create.EnableNuntioConnect = true
	create.DisableDefaultSignup = false
	create.DisableDefaultLogin = false
	create.ValidatePassword = true
	create.AuthConfig = &AuthConfig{
		WelcomeTitle:              "Welcome",
		WelcomeDetails:            "Already have an account? Welcome in! No? Create one below.",
		LoginButton:               "Login",
		LoginTitle:                "Login",
		LoginDetails:              "Fill in the details below to login to your account",
		RegisterButton:            "Register",
		RegisterTitle:             "Register",
		RegisterDetails:           "Fill in the details below to register for an account",
		MissingPasswordTitle:      "Missing required password",
		MissingPasswordDetails:    "You need to provide a password to create/login to an account.",
		MissingEmailTitle:         "Missing required email",
		MissingEmailDetails:       "You need to provide an email to create/login to an account",
		PasswordDoNotMatchTitle:   "Passwords do not match",
		PasswordDoNotMatchDetails: "The two passwords provided do not match.",
		CreatedBy:                 "Powered by Nuntio.",
	}
	if err := cr.EncryptConfig(actionCreate, create); err != nil {
		return nil, err
	}
	if _, err := cr.collection.InsertOne(ctx, create); err != nil {
		return nil, err
	}
	return config, nil
}
