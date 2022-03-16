package handler

import (
	"context"
	"github.com/softcorp-io/block-proto/go_block"
)

func (h *defaultHandler) UpdateSecurity(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error) {
	updatedUser, err := h.repository.UserRepository.UpdateSecurity(ctx, req.User, req.Update, req.EncryptionKey)
	if err != nil {
		return nil, err
	}
	return &go_block.UserResponse{
		User: updatedUser,
	}, nil
}
