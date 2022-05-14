package handler

import (
	"context"
	"github.com/nuntiodev/block-proto/go_block"
)

func (h *defaultHandler) GetText(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error) {
	text, err := h.repository.Text(ctx, req.Namespace)
	if err != nil {
		return nil, err
	}
	resp, err := text.Get(ctx, req.Text.Id)
	if err != nil {
		return nil, err
	}
	return &go_block.UserResponse{
		Text: resp,
	}, nil
}
