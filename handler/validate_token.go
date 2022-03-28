package handler

import (
	"context"
	"github.com/softcorp-io/block-proto/go_block"
	"github.com/softcorp-io/block-user-service/repository/token_repository"
)

func (h *defaultHandler) ValidateToken(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error) {
	customClaims, err := h.crypto.ValidateToken(req.Token.AccessToken)
	if err != nil {
		return nil, err
	}
	// validate if token is blocked in db
	tokens, err := h.repository.Tokens(ctx, req.Namespace)
	if err != nil {
		return nil, err
	}
	if err := tokens.IsBlocked(ctx, &token_repository.Token{
		AccessTokenId:  customClaims.Id,
		RefreshTokenId: customClaims.RefreshTokenId,
	}); err != nil {
		return nil, err
	}
	user := &go_block.User{
		Id: customClaims.UserId,
	}
	return &go_block.UserResponse{
		User: user,
	}, nil
}
