package handler

import (
	"context"
	"errors"
	"github.com/softcorp-io/block-proto/go_block"
	"github.com/softcorp-io/block-user-service/token"
	"golang.org/x/crypto/bcrypt"
)

func (h *defaultHandler) Login(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error) {
	resp, err := h.Get(ctx, req)
	if err != nil {
		return &go_block.UserResponse{}, err
	}
	if resp.User.Password == "" {
		return &go_block.UserResponse{}, errors.New("please update the user with a non-empty password")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(resp.User.Password), []byte(req.User.Password)); err != nil {
		return &go_block.UserResponse{}, err
	}
	// issue access and refresh token pair
	refreshToken, claims, err := h.token.GenerateToken(privateKey, resp.User.Id, "", token.TokenTypeRefresh, refreshTokenExpiry)
	if err != nil {
		return &go_block.UserResponse{}, err
	}
	accessToken, _, err := h.token.GenerateToken(privateKey, resp.User.Id, claims.UserId, token.TokenTypeAccess, accessTokenExpiry)
	if err != nil {
		return &go_block.UserResponse{}, err
	}
	return &go_block.UserResponse{
		Token: &go_block.Token{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
		User: resp.User,
	}, nil
}
