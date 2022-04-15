package handler

import (
	"context"
	"errors"

	"github.com/io-nuntio/block-proto/go_block"
	"github.com/nuntio-dev/nuntio-user-block/token"
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
		isBlocked, err := tokens.IsBlocked(ctx, &go_block.Token{
			Id:     customClaims.RefreshTokenId,
			UserId: customClaims.UserId,
		})
		if err != nil {
			return nil, err
		}
		if isBlocked {
			return nil, errors.New("token is blocked")
		}
	}
	// else we always validate if id of token is blocked
	isBlocked, err := tokens.IsBlocked(ctx, &go_block.Token{
		Id:     customClaims.Id,
		UserId: customClaims.UserId,
	})
	if err != nil {
		return nil, err
	}
	if isBlocked {
		return nil, errors.New("token is blocked")
	}
	if _, err := tokens.UpdateUsedAt(ctx, &go_block.Token{
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
