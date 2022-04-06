package handler

import (
	"context"
	"errors"
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
	// for access tokens we also validate if refresh token is blocked
	if customClaims.Type == token.TokenTypeAccess {
		isBlocked, err := tokens.IsBlocked(ctx, &token_repository.Token{
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
	isBlocked, err := tokens.IsBlocked(ctx, &token_repository.Token{
		Id:     customClaims.Id,
		UserId: customClaims.UserId,
	})
	if err != nil {
		return nil, err
	}
	if isBlocked {
		return nil, errors.New("token is blocked")
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
		newRefreshToken, newRefreshclaims, err := h.token.GenerateToken(privateKey, customClaims.UserId, "", token.TokenTypeRefresh, refreshTokenExpiry)
		if err != nil {
			return nil, err
		}
		refreshToken = newRefreshToken
		customClaims = newRefreshclaims
		// create refresh token in database
		if _, err := tokens.Create(ctx, &token_repository.Token{
			Id:        newRefreshclaims.Id,
			UserId:    newRefreshclaims.UserId,
			ExpiresAt: time.Unix(newRefreshclaims.ExpiresAt, 0),
		}); err != nil {
			return &go_block.UserResponse{}, err
		}
	}
	// generate new access token from refresh token
	newAccessToken, newAccessClaims, err := h.token.GenerateToken(privateKey, customClaims.UserId, customClaims.Id, token.TokenTypeAccess, accessTokenExpiry)
	if err != nil {
		return nil, err
	}
	// add new access token to database
	if _, err := tokens.Create(ctx, &token_repository.Token{
		Id:        newAccessClaims.Id,
		UserId:    newAccessClaims.UserId,
		ExpiresAt: time.Unix(newAccessClaims.ExpiresAt, 0),
	}); err != nil {
		return &go_block.UserResponse{}, err
	}
	if _, err := tokens.UpdateUsedAt(ctx, &token_repository.Token{
		Id: customClaims.Id,
	}); err != nil {
		return &go_block.UserResponse{}, err
	}
	return &go_block.UserResponse{
		Token: &go_block.Token{
			AccessToken:  newAccessToken,
			RefreshToken: refreshToken,
		},
	}, nil
}
