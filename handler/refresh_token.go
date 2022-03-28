package handler

import (
	"context"
	"github.com/softcorp-io/block-proto/go_block"
	"github.com/softcorp-io/block-user-service/crypto"
	"github.com/softcorp-io/block-user-service/repository/token_repository"
)

func (h *defaultHandler) RefreshToken(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error) {
	customClaims, err := h.crypto.ValidateToken(req.Token.RefreshToken)
	if err != nil {
		return nil, err
	}
	// validate if blocked in db
	tokens, err := h.repository.Tokens(ctx, req.Namespace)
	if err != nil {
		return nil, err
	}
	if err := tokens.IsBlocked(ctx, &token_repository.Token{
		Id: customClaims.Id,
	}); err != nil {
		return nil, err
	}
	// generate new access token from refresh token
	newAccessToken, err := h.crypto.GenerateToken(customClaims.UserId, crypto.TokenTypeAccess, h.accessTokenExpiry)
	if err != nil {
		return nil, err
	}
	return &go_block.UserResponse{
		Token: &go_block.Token{
			AccessToken:  newAccessToken,
			RefreshToken: req.Token.AccessToken,
		},
	}, nil
}
