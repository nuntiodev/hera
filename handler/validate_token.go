package handler

import (
	"context"
	"github.com/softcorp-io/block-proto/go_block"
	"github.com/softcorp-io/block-user-service/repository/token_repository"
	"github.com/softcorp-io/block-user-service/token"
)

func (h *defaultHandler) ValidateToken(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error) {
	customClaims, err := h.token.ValidateToken(publicKey, req.Token.AccessToken)
	if err != nil {
		return nil, err
	}
	// validate if token is blocked in db
	tokens, err := h.repository.Tokens(ctx, req.Namespace)
	if err != nil {
		return nil, err
	}
	// for access tokens we also validate if refresh token is blocked
	if customClaims.Type == token.TokenTypeAccess {
		if err := tokens.IsBlocked(ctx, &token_repository.Token{
			Id:     customClaims.RefreshTokenId,
			UserId: customClaims.UserId,
		}); err != nil {
			return nil, err
		}
	}
	// else we always validate if id of token is blocked
	if err := tokens.IsBlocked(ctx, &token_repository.Token{
		Id:     customClaims.Id,
		UserId: customClaims.UserId,
	}); err != nil {
		return nil, err
	}
	if _, err := tokens.UpdateUsedAt(ctx, &token_repository.Token{
		Id: customClaims.Id,
	}); err != nil {
		return &go_block.UserResponse{}, err
	}
	return &go_block.UserResponse{
		User: &go_block.User{
			Id: customClaims.UserId,
		},
	}, nil
}
