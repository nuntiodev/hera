package handler

import (
	"context"
	"github.com/nuntiodev/block-proto/go_block"
	"github.com/nuntiodev/nuntio-user-block/repository/token_repository"
)

/*
	GetTokens - this method returns information about all tokens for a specific user. UserId is included in the req.Token.
*/
func (h *defaultHandler) GetTokens(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error) {
	var (
		tokenRepo token_repository.TokenRepository
		tokens    []*go_block.Token
		err       error
	)
	tokenRepo, err = h.repository.Tokens(ctx, req.Namespace)
	if err != nil {
		return nil, err
	}
	tokens, err = tokenRepo.GetTokens(ctx, req.Token)
	return &go_block.UserResponse{
		Tokens: tokens,
	}, err
}
