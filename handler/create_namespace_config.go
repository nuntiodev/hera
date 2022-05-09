package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/nuntiodev/block-proto/go_block"
	"github.com/nuntiodev/nuntio-user-block/repository/email_repository"
)

func (h *defaultHandler) CreateNamespaceConfig(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error) {
	// create test user in namespace
	users, err := h.repository.Users().SetNamespace(req.Namespace).Build(ctx)
	if err != nil {
		return &go_block.UserResponse{}, fmt.Errorf("could not build user with err: %v", err)
	}
	metadata, err := json.Marshal(map[string]string{
		"role": "test",
	})
	if err != nil {
		return &go_block.UserResponse{}, fmt.Errorf("could not marshal user with err: %v", err)
	}
	// create test user
	if _, err := users.Create(ctx, &go_block.User{
		FirstName: "Test",
		LastName:  "User",
		Email:     "test@user.io",
		Metadata:  string(metadata),
	}); err != nil {
		return &go_block.UserResponse{}, err
	}
	// create initial config
	config, err := h.repository.Config(ctx, req.Namespace)
	if err != nil {
		return &go_block.UserResponse{}, fmt.Errorf("could not build config with err: %v", err)
	}
	createdConfig, err := config.Create(ctx, req.Config)
	if err != nil {
		return &go_block.UserResponse{}, fmt.Errorf("could not create config with err: %v", err)
	}
	// setup default emails
	emails, err := h.repository.Emails(ctx, req.Namespace)
	if err != nil {
		return &go_block.UserResponse{}, err
	}
	// create verification email
	if _, err := emails.Create(ctx, &go_block.Email{
		Id:             email_repository.VerificationEmail,
		Logo:           "",
		WelcomeMessage: "Hello",
		BodyMessage:    fmt.Sprintf("Thank you for signing up to %s. In order to get started, we ask of you to confirm your email by entering the following numbers in your %s app.", createdConfig.Name, createdConfig.Name),
		FooterMessage:  fmt.Sprintf("All the best from %s team", createdConfig.Name),
		Title:          "Verify your email",
		Subject:        "Verify your email account",
		TemplatePath:   emailVerificationTemplatePath,
	}); err != nil {
		return &go_block.UserResponse{}, err
	}
	// create verification email
	if _, err := emails.Create(ctx, &go_block.Email{
		Id:             email_repository.ResetPasswordEmail,
		Logo:           "",
		WelcomeMessage: "Hello",
		BodyMessage:    fmt.Sprintf("Thank you for signing up to %s. In order to reset your password, enter the following numbers in your %s app together with your new password.", createdConfig.Name, createdConfig.Name),
		FooterMessage:  fmt.Sprintf("All the best from %s team", createdConfig.Name),
		Title:          "Verify your email",
		Subject:        "Verify your email account",
		TemplatePath:   emailVerificationTemplatePath,
	}); err != nil {
		return &go_block.UserResponse{}, err
	}
	return &go_block.UserResponse{
		Config: createdConfig,
	}, nil
}
