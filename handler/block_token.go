package handler

import (
	"context"
	"errors"
	"github.com/softcorp-io/block-proto/go_block"
)

func (h *defaultHandler) BlockToken(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error) {
	token := ""
	if req.Token.AccessToken != "" {
		token = req.Token.AccessToken //todo: create new field just named token
	} else if req.Token.RefreshToken != "" {
		token = req.Token.RefreshToken
	} else {
		return &go_block.UserResponse{}, errors.New("no token in request")
	}
	customClaims, err := h.token.ValidateToken(publicKey, token)
	if err != nil {
		return &go_block.UserResponse{}, err
	}
	// validate if token is blocked in db
	tokens, err := h.repository.Tokens(ctx, req.Namespace)
	if err != nil {
		return &go_block.UserResponse{}, err
	}
	if _, err := tokens.Block(ctx, &go_block.Token{
		Id:     customClaims.Id,
		UserId: customClaims.UserId,
	}); err != nil {
		return &go_block.UserResponse{}, err
	}
	return &go_block.UserResponse{}, nil
}
