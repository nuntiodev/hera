package handler

import (
	"context"
	"github.com/nuntiodev/nuntio-user-block/repository/user_repository"

	"github.com/nuntiodev/block-proto/go_block"
)

/*
	UpdateUsername - this method is used to update the username of a user.
*/
func (h *defaultHandler) UpdateUsername(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error) {
	var (
		userRepo user_repository.UserRepository
		err      error
	)
	userRepo, err = h.repository.Users().SetNamespace(req.Namespace).Build(ctx)
	if err != nil {
		return &go_block.UserResponse{}, err
	}
	user, err := userRepo.UpdateUsername(ctx, req.User, req.Update)
	return &go_block.UserResponse{
		User: user,
	}, err
}
