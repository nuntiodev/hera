package handler

import (
	"context"

	"github.com/io-nuntio/block-proto/go_block"
)

func (h *defaultHandler) Heartbeat(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error) {
	if err := h.repository.Liveness(ctx); err != nil {
		return &go_block.UserResponse{}, err
	}
	return &go_block.UserResponse{}, nil
}
