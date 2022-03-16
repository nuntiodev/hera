package handler

import (
	"context"
	"github.com/softcorp-io/block-proto/go_block"
)

func (h *defaultHandler) GetAll(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error) {
	getUsers, err := h.repository.UserRepository.GetAll(ctx, req.Filter, req.Namespace, req.EncryptionKey)
	if err != nil {
		return nil, err
	}
	usersInNamespace, err := h.repository.UserRepository.Count(ctx, req.Namespace)
	if err != nil {
		return nil, err
	}
	return &go_block.UserResponse{
		Users:       getUsers,
		UsersAmount: usersInNamespace,
	}, nil
}
