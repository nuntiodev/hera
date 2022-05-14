package handler

import (
	"context"
	"errors"
	"fmt"

	"github.com/nuntiodev/block-proto/go_block"
)

func (h *defaultHandler) UpdateEmail(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error) {
	// get config
	config, err := h.repository.Config(ctx, req.Namespace)
	if err != nil {
		return &go_block.UserResponse{}, fmt.Errorf("could not build config repository with err: %v", err)
	}
	namespaceConfig, err := config.GetNamespaceConfig(ctx)
	if err != nil {
		return &go_block.UserResponse{}, fmt.Errorf("could not get namespace config with err: %v", err)
	}
	// we cannot send an email if the email provider is not enabled
	if !h.emailEnabled && namespaceConfig.RequireEmailVerification {
		return &go_block.UserResponse{}, errors.New("email provider is not enabled and verification email cannot be sent")
	} else if namespaceConfig.RequireEmailVerification && req.Update.Email == "" {
		return &go_block.UserResponse{}, errors.New("require email is enabled and email is empty")
	}
	users, err := h.repository.Users().SetNamespace(req.Namespace).SetEncryptionKey(req.EncryptionKey).Build(ctx)
	if err != nil {
		return &go_block.UserResponse{}, fmt.Errorf("could not build user repository with err: %v", err)
	}
	// perform update
	updatedUser, err := users.UpdateEmail(ctx, req.User, req.Update)
	if err != nil {
		return &go_block.UserResponse{}, fmt.Errorf("could not update user email with err: %v", err)
	}
	return &go_block.UserResponse{
		User: updatedUser,
	}, nil
}
