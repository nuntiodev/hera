package handler

import (
	"context"
	"github.com/softcorp-io/block-proto/go_block"
)

func (h *defaultHandler) Get(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error) {
	users, err := h.repository.Users(ctx, req.Namespace, req.EncryptionKey)
	if err != nil {
		return &go_block.UserResponse{}, err
	}
	getUser, err := users.Get(ctx, req.User)
	if err != nil {
		return nil, err
	}
	return &go_block.UserResponse{
		User: getUser,
	}, nil
}
