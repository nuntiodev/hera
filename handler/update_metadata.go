package handler

import (
	"context"
	"github.com/nuntiodev/nuntio-user-block/repository/user_repository"

	"github.com/nuntiodev/block-proto/go_block"
)

/*
	UpdateMetadata - this method updates a users metadata which is stored as JSON and used to store additional
	information about a user.
*/
func (h *defaultHandler) UpdateMetadata(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error) {
	var (
		userRepo user_repository.UserRepository
		user     *go_block.User
		err      error
	)
	userRepo, err = h.repository.Users().SetNamespace(req.Namespace).SetEncryptionKey(req.EncryptionKey).Build(ctx)
	if err != nil {
		return &go_block.UserResponse{}, err
	}
	user, err = userRepo.UpdateMetadata(ctx, req.User, req.Update)
	return &go_block.UserResponse{
		User: user,
	}, err
}
