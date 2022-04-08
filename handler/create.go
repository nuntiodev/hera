package handler

import (
	"context"
	"github.com/softcorp-io/block-proto/go_block"
)

func (h *defaultHandler) Create(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error) {
	users, err := h.repository.Users().SetNamespace(req.Namespace).SetEncryptionKey(req.EncryptionKey).WithPasswordValidation(req.ValidatePassword || validatePassword).Build(ctx)
	if err != nil {
		return &go_block.UserResponse{}, err
	}
	createdUser, err := users.Create(ctx, req.User)
	if err != nil {
		return &go_block.UserResponse{}, err
	}
	return &go_block.UserResponse{
		User: createdUser,
	}, nil
}
