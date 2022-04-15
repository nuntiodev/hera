package handler

import (
	"context"
	"errors"
	"time"

	"github.com/nuntiodev/block-proto/go_block"
	"github.com/nuntiodev/nuntio-user-block/token"
	"golang.org/x/crypto/bcrypt"
	ts "google.golang.org/protobuf/types/known/timestamppb"
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
	refreshToken, refreshClaims, err := h.token.GenerateToken(privateKey, resp.User.Id, "", token.TokenTypeRefresh, refreshTokenExpiry)
	if err != nil {
		return &go_block.UserResponse{}, err
	}
	accessToken, accessClaims, err := h.token.GenerateToken(privateKey, resp.User.Id, refreshClaims.UserId, token.TokenTypeAccess, accessTokenExpiry)
	if err != nil {
		return &go_block.UserResponse{}, err
	}
	// setup token database
	tokens, err := h.repository.Tokens(ctx, req.Namespace)
	if err != nil {
		return &go_block.UserResponse{}, err
	}
	// create refresh token in database
	if _, err := tokens.Create(ctx, &go_block.Token{
		Id:        refreshClaims.Id,
		UserId:    refreshClaims.UserId,
		ExpiresAt: ts.New(time.Unix(refreshClaims.ExpiresAt, 0)),
	}); err != nil {
		return &go_block.UserResponse{}, err
	}
	// create access token in database
	if _, err := tokens.Create(ctx, &go_block.Token{
		Id:        accessClaims.Id,
		UserId:    accessClaims.UserId,
		ExpiresAt: ts.New(time.Unix(accessClaims.ExpiresAt, 0)),
	}); err != nil {
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
