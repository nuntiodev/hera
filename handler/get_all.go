package handler

import (
	"context"
	"github.com/softcorp-io/block-proto/go_block"
)

func (h *defaultHandler) GetAll(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error) {
	users, err := h.repository.Users().SetNamespace(req.Namespace).SetEncryptionKey(req.EncryptionKey).Build(ctx)
	if err != nil {
		return &go_block.UserResponse{}, err
	}
	getUsers, err := users.GetAll(ctx, req.Filter)
	if err != nil {
		return nil, err
	}
	usersInNamespace, err := users.Count(ctx)
	if err != nil {
		return nil, err
	}
	return &go_block.UserResponse{
		Users:       getUsers,
		UsersAmount: usersInNamespace,
	}, nil
}
