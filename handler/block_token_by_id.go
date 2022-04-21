package handler

import (
	"context"

	"github.com/nuntiodev/block-proto/go_block"
)

func (h *defaultHandler) BlockTokenById(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error) {
	tokens, err := h.repository.Tokens(ctx, req.Namespace)
	if err != nil {
		return &go_block.UserResponse{}, err
	}
	if _, err := tokens.Block(ctx, &go_block.Token{
		Id:     req.Token.Id,
		UserId: req.Token.UserId,
	}); err != nil {
		return &go_block.UserResponse{}, err
	}
	return &go_block.UserResponse{}, nil
}
