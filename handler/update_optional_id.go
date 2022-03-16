package handler

import (
	"context"
	"github.com/softcorp-io/block-proto/go_block"
)

func (h *defaultHandler) UpdateOptionalId(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error) {
	updatedUser, err := h.repository.UserRepository.UpdateOptionalId(ctx, req.User, req.Update)
	if err != nil {
		return nil, err
	}
	return &go_block.UserResponse{
		User: updatedUser,
	}, nil
}
