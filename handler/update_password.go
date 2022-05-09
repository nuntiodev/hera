package handler

import (
	"context"
	"github.com/nuntiodev/block-proto/go_block"
)

func (h *defaultHandler) UpdatePassword(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error) {
	// get config
	config, err := h.repository.Config(ctx, req.Namespace)
	if err != nil {
		return &go_block.UserResponse{}, err
	}
	namespaceConfig, err := config.GetNamespaceConfig(ctx)
	if err != nil {
		return &go_block.UserResponse{}, err
	}
	users, err := h.repository.Users().SetNamespace(req.Namespace).WithPasswordValidation(namespaceConfig.ValidatePassword).Build(ctx)
	if err != nil {
		return &go_block.UserResponse{}, err
	}
	updatedUser, err := users.UpdatePassword(ctx, req.User, req.Update)
	if err != nil {
		return &go_block.UserResponse{}, err
	}
	return &go_block.UserResponse{
		User: updatedUser,
	}, nil
}
