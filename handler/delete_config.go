package handler

import (
	"context"
	"github.com/nuntiodev/block-proto/go_block"
)

func (h *defaultHandler) DeleteConfig(ctx context.Context, req *go_block.ConfigRequest) (*go_block.ConfigResponse, error) {
	config, err := h.repository.Config(ctx, req.Namespace)
	if err != nil {
		return nil, err
	}
	if err := config.Delete(ctx, req.Config); err != nil {
		return nil, err
	}
	return &go_block.ConfigResponse{}, nil
}
