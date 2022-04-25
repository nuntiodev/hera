package handler

import (
	"context"
	"github.com/nuntiodev/block-proto/go_block"
)

func (h *defaultHandler) UpdateConfigSettings(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error) {
	config, err := h.repository.Config(ctx, req.Namespace)
	if err != nil {
		return nil, err
	}
	resp, err := config.UpdateSettings(ctx, req.Config)
	if err != nil {
		return nil, err
	}
	return &go_block.UserResponse{
		Config: resp,
	}, nil
}
