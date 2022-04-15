package handler

import (
	"context"

	"github.com/io-nuntio/block-proto/go_block"
)

func (h *defaultHandler) UpdateSecurity(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error) {
	users, err := h.repository.Users().SetNamespace(req.Namespace).SetEncryptionKey(req.EncryptionKey).Build(ctx)
	if err != nil {
		return &go_block.UserResponse{}, err
	}
	updatedUser, err := users.UpdateSecurity(ctx, req.User)
	if err != nil {
		return nil, err
	}
	return &go_block.UserResponse{
		User: updatedUser,
	}, nil
}
