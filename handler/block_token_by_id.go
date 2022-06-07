package handler

import (
	"context"
	"github.com/nuntiodev/nuntio-user-block/repository/token_repository"

	"github.com/nuntiodev/block-proto/go_block"
)

/*
	BlockTokenById - this method will block an access or refresh token in the database for a specific login-session.
*/
func (h *defaultHandler) BlockTokenById(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error) {
	var (
		tokenRepo token_repository.TokenRepository
		err       error
	)
	tokenRepo, err = h.repository.Tokens(ctx, req.Namespace, req.EncryptionKey)
	if err != nil {
		return &go_block.UserResponse{}, err
	}
	if _, err = tokenRepo.Block(ctx, &go_block.Token{
		Id:     req.Token.Id,
		UserId: req.Token.UserId,
	}); err != nil {
		return &go_block.UserResponse{}, err
	}
	return &go_block.UserResponse{}, nil
}
