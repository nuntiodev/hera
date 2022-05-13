package handler

import (
	"context"
	"errors"
	"github.com/nuntiodev/block-proto/go_block"
)

func (h *defaultHandler) Create(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error) {
	// get config
	config, err := h.repository.Config(ctx, req.Namespace)
	if err != nil {
		return &go_block.UserResponse{}, err
	}
	namespaceConfig, err := config.GetNamespaceConfig(ctx)
	if err != nil {
		return &go_block.UserResponse{}, err
	}
	// we cannot send an email if the email provider is not enabled
	if !h.emailEnabled && namespaceConfig.RequireEmailVerification {
		return &go_block.UserResponse{}, errors.New("email provider is not enabled and verification email cannot be sent")
	} else if namespaceConfig.RequireEmailVerification && req.User.Email == "" {
		return &go_block.UserResponse{}, errors.New("require email verification is enabled and user email is empty")
	}
	// create user in db
	users, err := h.repository.Users().SetNamespace(req.Namespace).SetEncryptionKey(req.EncryptionKey).WithPasswordValidation(namespaceConfig.ValidatePassword).Build(ctx)
	if err != nil {
		return &go_block.UserResponse{}, err
	}
	createdUser, err := users.Create(ctx, req.User)
	if err != nil {
		return &go_block.UserResponse{}, err
	}
	if h.emailEnabled && namespaceConfig.RequireEmailVerification { // email is enabled, and we require email verification
		req.User.Id = createdUser.Id
		if _, err := h.SendVerificationEmail(ctx, req); err != nil {
			return nil, err
		}
	}
	return &go_block.UserResponse{
		User: createdUser,
	}, nil
}
