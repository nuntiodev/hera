package handler

import (
	"context"
	"fmt"

	"github.com/nuntiodev/block-proto/go_block"
)

func (h *defaultHandler) DeleteNamespace(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error) {
	users, err := h.repository.Users().SetNamespace(req.Namespace).Build(ctx)
	if err != nil {
		return &go_block.UserResponse{}, fmt.Errorf("could not build user with err: %v", err)
	}
	// also delete config
	config, err := h.repository.Config(ctx, req.Namespace)
	if err != nil {
		return &go_block.UserResponse{}, fmt.Errorf("could not build config with err: %v", err)
	}
	if err := config.Delete(ctx); err != nil {
		return &go_block.UserResponse{}, fmt.Errorf("could not delete config with err: %v", err)
	}
	if err := users.DeleteAll(ctx); err != nil {
		return &go_block.UserResponse{}, fmt.Errorf("could not delete namespace with err: %v", err)
	}
	return &go_block.UserResponse{}, nil
}
