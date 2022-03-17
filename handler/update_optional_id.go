package handler

import (
	"context"
	"github.com/softcorp-io/block-proto/go_block"
)

func (h *defaultHandler) UpdateOptionalId(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error) {
	users, err := h.repository.Users(ctx, req.Namespace)
	if err != nil {
		return &go_block.UserResponse{}, err
	}
	updatedUser, err := users.UpdateOptionalId(ctx, req.User, req.Update)
	if err != nil {
		return nil, err
	}
	return &go_block.UserResponse{
		User: updatedUser,
	}, nil
}
