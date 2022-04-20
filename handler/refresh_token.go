package handler

import (
	"context"
	"errors"
	"time"

	"github.com/nuntiodev/block-proto/go_block"
	"github.com/nuntiodev/nuntio-user-block/token"
	ts "google.golang.org/protobuf/types/known/timestamppb"
)

func (h *defaultHandler) RefreshToken(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error) {
	refreshClaims, err := h.token.ValidateToken(publicKey, req.Token.RefreshToken)
	if err != nil {
		return nil, err
	}
	// validate if blocked in db
	tokens, err := h.repository.Tokens(ctx, req.Namespace)
	if err != nil {
		return nil, err
	}
	// we can only use a refresh token to generate a new one
	if refreshClaims.Type != token.TokenTypeRefresh {
		return nil, errors.New("invalid refresh token")
	}
	// else we always validate if id of token is blocked
	isBlocked, err := tokens.IsBlocked(ctx, &go_block.Token{
		Id:     refreshClaims.Id,
		UserId: refreshClaims.UserId,
	})
	if err != nil {
		return nil, err
	}
	if isBlocked {
		return nil, errors.New("token is blocked")
	}
	// build data for token
	loggedInFrom := ""
	deviceInfo := ""
	if req.Token != nil {
		loggedInFrom = req.Token.LoggedInFrom
		deviceInfo = req.Token.DeviceInfo
	}
	// if refresh token is about to expire (in less than 10 hours), create a new one and block the old one
	refreshToken := req.Token.RefreshToken
	if time.Unix(refreshClaims.ExpiresAt, 0).Sub(time.Now()) < time.Hour*10 {
		if _, err := h.BlockToken(ctx, &go_block.UserRequest{
			Token: &go_block.Token{
				RefreshToken: refreshToken,
			},
		}); err != nil {
			return nil, err
		}
		newRefreshToken, newRefreshclaims, err := h.token.GenerateToken(privateKey, refreshClaims.UserId, "", token.TokenTypeRefresh, refreshTokenExpiry)
		if err != nil {
			return nil, err
		}
		refreshToken = newRefreshToken
		refreshClaims = newRefreshclaims
		// create refresh token in database
		if _, err := tokens.Create(ctx, &go_block.Token{
			Id:           newRefreshclaims.Id,
			UserId:       newRefreshclaims.UserId,
			Type:         go_block.TokenType_TOKEN_TYPE_REFRESH,
			LoggedInFrom: loggedInFrom,
			DeviceInfo:   deviceInfo,
			ExpiresAt:    ts.New(time.Unix(newRefreshclaims.ExpiresAt, 0)),
		}); err != nil {
			return &go_block.UserResponse{}, err
		}
	}
	// generate new access token from refresh token
	newAccessToken, newAccessClaims, err := h.token.GenerateToken(privateKey, refreshClaims.UserId, refreshClaims.Id, token.TokenTypeAccess, accessTokenExpiry)
	if err != nil {
		return nil, err
	}
	// add new access token to database
	if _, err := tokens.Create(ctx, &go_block.Token{
		Id:           newAccessClaims.Id,
		UserId:       newAccessClaims.UserId,
		Type:         go_block.TokenType_TOKEN_TYPE_ACCESS,
		LoggedInFrom: loggedInFrom,
		DeviceInfo:   deviceInfo,
		ExpiresAt:    ts.New(time.Unix(newAccessClaims.ExpiresAt, 0)),
	}); err != nil {
		return &go_block.UserResponse{}, err
	}
	// set refresh token used at
	if _, err := tokens.UpdateUsedAt(ctx, &go_block.Token{
		Id:     refreshClaims.Id,
		UserId: refreshClaims.UserId,
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
