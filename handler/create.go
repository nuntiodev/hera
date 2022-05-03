package handler

import (
	"context"
	"errors"
	"github.com/nuntiodev/nuntio-user-block/email"
	"strings"

	"github.com/nuntiodev/block-proto/go_block"
)

func (h *defaultHandler) Create(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error) {
	// we cannot send an email if the email provider is not enabled
	if !h.emailEnabled && req.RequireEmailVerification {
		return &go_block.UserResponse{}, errors.New("email provider is not enabled and verification email cannot be sent")
	}
	users, err := h.repository.Users().SetNamespace(req.Namespace).SetEncryptionKey(req.EncryptionKey).WithPasswordValidation(req.ValidatePassword || validatePassword).Build(ctx)
	if err != nil {
		return &go_block.UserResponse{}, err
	}
	createdUser, err := users.Create(ctx, req.User)
	if err != nil {
		return &go_block.UserResponse{}, err
	}
	if h.emailEnabled && req.RequireEmailVerification { // email is enabled, and we require email verification
		//todo: provide subject of creation emails in email repository config
		nameOfUser := createdUser.Email
		if createdUser.FirstName != "" {
			nameOfUser = strings.TrimSpace(createdUser.FirstName + " " + createdUser.LastName)
		}
		if err := h.email.SendEmail(createdUser.Email, "", "", &email.TemplateData{
			LogoUrl:        "",
			WelcomeMessage: "",
			NameOfUser:     nameOfUser,
			BodyMessage:    "",
			FooterMessage:  "",
		}); err != nil {
			if err != nil {
				return &go_block.UserResponse{}, err
			}
		}
		if _, err := users.UpdateVerificationEmailSent(ctx, createdUser); err != nil {
			return &go_block.UserResponse{}, err
		}
	}
	return &go_block.UserResponse{
		User: createdUser,
	}, nil
}
