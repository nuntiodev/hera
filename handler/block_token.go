package handler

import (
	"context"
	"github.com/softcorp-io/block-proto/go_block"
	"github.com/softcorp-io/block-user-service/repository/token_repository"
)

func (h *defaultHandler) BlockToken(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error) {
	customClaims, err := h.crypto.ValidateToken(req.Token.AccessToken)
	if err != nil {
		return nil, err
	}
	// validate if token is blocked in db
	tokens, err := h.repository.Tokens(ctx, req.Namespace)
	if err != nil {
		return nil, err
	}
	if err := tokens.BlockToken(ctx, &token_repository.Token{
		Id:        customClaims.Id,
		ExpiresAt: customClaims.ExpiresAt,
	}); err != nil {
		return nil, err
	}
	return &go_block.UserResponse{}, nil
}
