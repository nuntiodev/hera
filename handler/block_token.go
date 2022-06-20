package handler

import (
	"context"
	"fmt"
	"github.com/nuntiodev/block-proto/go_block"
	"github.com/nuntiodev/nuntio-user-block/repository/token_repository"
)

/*
	BlockToken - this method will block an access or refresh token in the database for a specific login-session.
*/
func (h *defaultHandler) BlockToken(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error) {
	var (
		claims    *go_block.CustomClaims
		tokenRepo token_repository.TokenRepository
		err       error
	)
	// validate requested token and get id of the token
	claims, err = h.token.ValidateToken(publicKey, req.TokenPointer)
	if err != nil {
		return &go_block.UserResponse{}, err
	}
	// validate if token is blocked in db
	fmt.Println(req.TokenPointer)
	fmt.Println(claims)
	tokenRepo, err = h.repository.Tokens(ctx, req.Namespace, req.EncryptionKey)
	if err != nil {
		return &go_block.UserResponse{}, err
	}
	_, err = tokenRepo.Block(ctx, &go_block.Token{
		Id:     claims.Id,
		UserId: claims.UserId,
	})
	return &go_block.UserResponse{}, err
}
