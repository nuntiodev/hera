package handler

import (
	"context"
	"errors"
	"github.com/nuntiodev/nuntio-user-block/repository/user_repository"

	"github.com/nuntiodev/block-proto/go_block"
	"golang.org/x/crypto/bcrypt"
)

/*
	ValidateCredentials - this method is used to validate the credentials of a user without issuing a token.
*/
func (h *defaultHandler) ValidateCredentials(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error) {
	var (
		userRepo user_repository.UserRepository
		user     *go_block.User
		err      error
	)
	userRepo, err = h.repository.Users().SetNamespace(req.Namespace).SetEncryptionKey(req.EncryptionKey).Build(ctx)
	if err != nil {
		return &go_block.UserResponse{}, err
	}
	user, err = userRepo.Get(ctx, req.User, true)
	if err != nil {
		return &go_block.UserResponse{}, err
	}
	if user.Password == "" {
		return &go_block.UserResponse{}, errors.New("please update the user with a non-empty password")
	}
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.User.Password)); err != nil {
		return &go_block.UserResponse{}, err
	}
	return &go_block.UserResponse{
		User: user,
	}, nil
}
