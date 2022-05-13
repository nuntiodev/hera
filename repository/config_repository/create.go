package config_repository

import (
	"context"
	"errors"
	"github.com/nuntiodev/block-proto/go_block"
	ts "google.golang.org/protobuf/types/known/timestamppb"
)

func (c *defaultConfigRepository) Create(ctx context.Context, config *go_block.Config) (*go_block.Config, error) {
	prepare(actionCreate, config)
	if config == nil {
		return nil, errors.New("missing required config")
	} else if config.Id == "" {
		return nil, errors.New("missing required config")
	}
	// set default fields
	config.EnableNuntioConnect = true
	config.DisableDefaultSignup = false
	config.DisableDefaultLogin = false
	config.ValidatePassword = true
	config.RequireEmailVerification = true
	config.CreatedAt = ts.Now()
	config.UpdatedAt = ts.Now()
	// set default text
	config.GeneralText = &go_block.GeneralText{
		MissingPasswordTitle:   "Missing required password",
		MissingPasswordDetails: "You need to provide a password to create/login to an account.",
		MissingEmailTitle:      "Missing required email",
		MissingEmailDetails:    "You need to provide an email to create/login to an account",
		CreatedBy:              "Powered by Nuntio.",
		PasswordHint:           "Enter your password",
		EmailHint:              "Enter your email",
		ErrorTitle:             "An error occurred",
		ErrorDescription:       "Something went wrong. Please try again.",
		NoWifiTitle:            "No connection",
		NoWifiDescription:      "No stable cellular or wifi connection is available.",
	}
	config.WelcomeText = &go_block.WelcomeText{
		WelcomeTitle:       "Welcome",
		WelcomeDetails:     "Welcome to this awesome platform.",
		ContinueWithNuntio: "Continue with",
	}
	config.RegisterText = &go_block.RegisterText{
		RegisterButton:            "Register",
		RegisterTitle:             "Register",
		RegisterDetails:           "Fill in the fields in order to register for an account",
		PasswordDoNotMatchTitle:   "Passwords do not match",
		PasswordDoNotMatchDetails: "The provided and repeat passwords are not the same.",
		RepeatPasswordHint:        "Enter your password again",
		ContainsSpecialChar:       "Password must contain a special char",
		ContainsNumberChar:        "Password must contain a number",
		PasswordMustMatch:         "The two passwords must match",
		ContainsEightChars:        "Password must be at least 8 chars long",
	}
	config.LoginText = &go_block.LoginText{
		LoginButton:    "Login",
		LoginTitle:     "Login",
		LoginDetails:   "Fill in the details below to login to your account",
		ForgotPassword: "Forgot your password?",
	}
	create := ProtoConfigToConfig(config)
	if len(c.internalEncryptionKeys) > 0 {
		if err := c.EncryptConfig(actionCreate, create); err != nil {
			return nil, err
		}
		create.InternalEncryptionLevel = int32(len(c.internalEncryptionKeys))
	}
	if _, err := c.collection.InsertOne(ctx, create); err != nil {
		return nil, err
	}
	// set created fields
	config.InternalEncryptionLevel = create.InternalEncryptionLevel
	return config, nil
}
