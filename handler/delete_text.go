package handler

import (
	"context"
	"github.com/nuntiodev/block-proto/go_block"
)

func (h *defaultHandler) DeleteText(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error) {
	text, err := h.repository.Text(ctx, req.Namespace)
	if err != nil {
		return nil, err
	}
	if err := text.Delete(ctx, req.Text.Id); err != nil {
		return nil, err
	}
	return &go_block.UserResponse{}, nil
}
