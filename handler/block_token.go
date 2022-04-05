package handler

import (
	"context"
	"errors"
	"github.com/softcorp-io/block-proto/go_block"
	"github.com/softcorp-io/block-user-service/repository/token_repository"
	"github.com/softcorp-io/x/cryptox"
)

func (h *defaultHandler) BlockToken(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error) {
	token := ""
	if req.Token.AccessToken != "" {
		token = req.Token.AccessToken
	} else if req.Token.RefreshToken != "" {
		token = req.Token.RefreshToken
	} else {
		return &go_block.UserResponse{}, errors.New("no token in request")
	}
	customClaims, err := h.token.ValidateToken(publicKey, token)
	if err != nil {
		return &go_block.UserResponse{}, err
	}
	// build ids
	accessTokenId := ""
	refreshTokenId := ""
	if customClaims.Type == cryptox.TokenTypeAccess {
		accessTokenId = customClaims.Id
	} else if customClaims.Type == cryptox.TokenTypeRefresh {
		refreshTokenId = customClaims.Id
	}
	// validate if token is blocked in db
	tokens, err := h.repository.Tokens(ctx, req.Namespace)
	if err != nil {
		return &go_block.UserResponse{}, err
	}
	if err := tokens.BlockToken(ctx, &token_repository.Token{
		AccessTokenId:  accessTokenId,
		RefreshTokenId: refreshTokenId,
		ExpiresAt:      customClaims.ExpiresAt,
	}); err != nil {
		return &go_block.UserResponse{}, err
	}
	return &go_block.UserResponse{}, nil
}
