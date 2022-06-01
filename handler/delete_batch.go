package handler

import (
	"context"
	"github.com/nuntiodev/nuntio-user-block/repository/user_repository"

	"github.com/nuntiodev/block-proto/go_block"
)

/*
	DeleteBatch - this method deletes a batch of users.
*/
func (h *defaultHandler) DeleteBatch(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error) {
	var (
		userRepo user_repository.UserRepository
		err      error
	)
	userRepo, err = h.repository.Users().SetNamespace(req.Namespace).Build(ctx)
	if err != nil {
		return &go_block.UserResponse{}, err
	}
	return &go_block.UserResponse{}, userRepo.DeleteBatch(ctx, req.UserBatch)
}
