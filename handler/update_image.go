package handler

import (
	"context"
	"github.com/nuntiodev/nuntio-user-block/repository/user_repository"

	"github.com/nuntiodev/block-proto/go_block"
)

/*
	UpdateImage - this method updates a users image.
*/
func (h *defaultHandler) UpdateImage(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error) {
	var (
		userRepo user_repository.UserRepository
		user     *go_block.User
		err      error
	)
	userRepo, err = h.repository.Users().SetNamespace(req.Namespace).SetEncryptionKey(req.EncryptionKey).Build(ctx)
	if err != nil {
		return &go_block.UserResponse{}, err
	}
	user, err = userRepo.UpdateImage(ctx, req.User, req.Update)
	return &go_block.UserResponse{
		User: user,
	}, err
}
