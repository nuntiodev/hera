package handler

import (
	"context"
	"github.com/softcorp-io/block-proto/go_block"
	"github.com/softcorp-io/block-user-service/repository/token_repository"
	"github.com/softcorp-io/block-user-service/token"
	"time"
)

func (h *defaultHandler) RefreshToken(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error) {
	customClaims, err := h.token.ValidateToken(publicKey, req.Token.RefreshToken)
	if err != nil {
		return nil, err
	}
	// validate if blocked in db
	tokens, err := h.repository.Tokens(ctx, req.Namespace)
	if err != nil {
		return nil, err
	}
	if err := tokens.IsBlocked(ctx, &token_repository.Token{
		RefreshTokenId: customClaims.Id,
	}); err != nil {
		return nil, err
	}
	// if refresh token is about to expire (in less than 10 hours), create a new one and block the old one
	refreshToken := req.Token.RefreshToken
	if time.Unix(customClaims.ExpiresAt, 0).Sub(time.Now()) < time.Hour*10 {
		if _, err := h.BlockToken(ctx, &go_block.UserRequest{
			Token: &go_block.Token{
				RefreshToken: refreshToken,
			},
		}); err != nil {
			return nil, err
		}
		newRefreshToken, claims, err := h.token.GenerateToken(privateKey, customClaims.UserId, "", token.TokenTypeRefresh, refreshTokenExpiry)
		if err != nil {
			return nil, err
		}
		refreshToken = newRefreshToken
		customClaims = claims
	}
	// generate new access token from refresh token
	newAccessToken, _, err := h.token.GenerateToken(privateKey, customClaims.UserId, customClaims.Id, token.TokenTypeAccess, accessTokenExpiry)
	if err != nil {
		return nil, err
	}
	return &go_block.UserResponse{
		Token: &go_block.Token{
			AccessToken:  newAccessToken,
			RefreshToken: refreshToken,
		},
	}, nil
}
