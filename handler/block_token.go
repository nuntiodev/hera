package handler

import (
	"context"

	"github.com/nuntiodev/block-proto/go_block"
)

func (h *defaultHandler) BlockToken(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error) {
	customClaims, err := h.token.ValidateToken(publicKey, req.TokenPointer)
	if err != nil {
		return &go_block.UserResponse{}, err
	}
	// validate if token is blocked in db
	tokens, err := h.repository.Tokens(ctx, req.Namespace)
	if err != nil {
		return &go_block.UserResponse{}, err
	}
	if _, err := tokens.Block(ctx, &go_block.Token{
		Id:     customClaims.Id,
		UserId: customClaims.UserId,
	}); err != nil {
		return &go_block.UserResponse{}, err
	}
	return &go_block.UserResponse{}, nil
}
