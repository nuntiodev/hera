package handler

import (
	"context"

	"github.com/nuntiodev/block-proto/go_block"
)

func (h *defaultHandler) GetTokens(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error) {
	tokens, err := h.repository.Tokens(ctx, req.Namespace)
	if err != nil {
		return nil, err
	}
	getTokens, err := tokens.GetTokens(ctx, &go_block.Token{
		UserId: req.User.Id,
	})
	if err != nil {
		return nil, err
	}
	return &go_block.UserResponse{
		Tokens: getTokens,
	}, nil
}
