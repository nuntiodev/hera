package handler

import (
	"context"
	"errors"

	"github.com/nuntiodev/block-proto/go_block"
	"github.com/nuntiodev/nuntio-user-block/token"
)

func (h *defaultHandler) ValidateToken(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error) {
	customClaims, err := h.token.ValidateToken(publicKey, req.TokenPointer)
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
	// build data for token
	loggedInFrom := &go_block.Location{}
	deviceInfo := ""
	if req.Token != nil {
		loggedInFrom = req.Token.LoggedInFrom
		deviceInfo = req.Token.DeviceInfo
	}
	if err != nil {
		return nil, err
	}
	if isBlocked {
		return nil, errors.New("token is blocked")
	}
	if _, err := tokens.UpdateUsedAt(ctx, &go_block.Token{
		Id:           customClaims.Id,
		LoggedInFrom: loggedInFrom,
		DeviceInfo:   deviceInfo,
	}); err != nil {
		return &go_block.UserResponse{}, err
	}
	// validate user exists
	get, err := h.Get(ctx, &go_block.UserRequest{
		User: &go_block.User{
			Id: customClaims.UserId,
		},
	})
	if err != nil {
		return nil, err
	}
	return get, nil
}
