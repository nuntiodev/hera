package handler

import (
	"context"
	"github.com/softcorp-io/block-proto/go_block"
)

func (h *defaultHandler) Get(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error) {
	getUser, err := h.repository.UserRepository.Get(ctx, req.User, req.EncryptionKey)
	if err != nil {
		return nil, err
	}
	return &go_block.UserResponse{
		User: getUser,
	}, nil
}
