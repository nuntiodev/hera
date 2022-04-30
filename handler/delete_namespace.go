package handler

import (
	"context"

	"github.com/nuntiodev/block-proto/go_block"
)

func (h *defaultHandler) DeleteNamespace(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error) {
	users, err := h.repository.Users().SetNamespace(req.Namespace).Build(ctx)
	if err != nil {
		return &go_block.UserResponse{}, err
	}
	// also delete config
	config, err := h.repository.Config(ctx, req.Namespace)
	if err != nil {
		return &go_block.UserResponse{}, err
	}
	if err := config.Delete(ctx, &go_block.Config{
		Id: req.Namespace,
	}); err != nil {
		return &go_block.UserResponse{}, err
	}
	return &go_block.UserResponse{}, users.DeleteAll(ctx)
}
