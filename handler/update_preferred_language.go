package handler

import (
	"context"
	"github.com/nuntiodev/nuntio-user-block/models"
	"github.com/nuntiodev/nuntio-user-block/repository/user_repository"

	"github.com/nuntiodev/block-proto/go_block"
)

/*
	UpdatePreferredLanguage - this method updates a users preferred language.
*/
func (h *defaultHandler) UpdatePreferredLanguage(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error) {
	var (
		userRepo user_repository.UserRepository
		user     *models.User
		err      error
	)
	userRepo, err = h.repository.UserRepositoryBuilder().SetNamespace(req.Namespace).Build(ctx)
	if err != nil {
		return &go_block.UserResponse{}, err
	}
	user, err = userRepo.UpdatePreferredLanguage(ctx, req.User, req.Update)
	if err != nil {
		return nil, err
	}
	return &go_block.UserResponse{
		User: models.UserToProtoUser(user),
	}, nil
}
