package handler

import (
	"context"
	"github.com/nuntiodev/block-proto/go_block"
	"github.com/nuntiodev/nuntio-user-block/models"
	"github.com/nuntiodev/nuntio-user-block/repository/user_repository"
)

/*
	UpdatePhoneNumber - this method updates a users phone number.
*/
func (h *defaultHandler) UpdatePhoneNumber(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error) {
	var (
		userRepo user_repository.UserRepository
		user     *models.User
		err      error
	)
	userRepo, err = h.repository.UserRepositoryBuilder().SetNamespace(req.Namespace).SetEncryptionKey(req.EncryptionKey).Build(ctx)
	if err != nil {
		return &go_block.UserResponse{}, err
	}
	// perform update
	user, err = userRepo.UpdatePhoneNumber(ctx, req.User, req.Update)
	if err != nil {
		return nil, err
	}
	return &go_block.UserResponse{
		User: models.UserToProtoUser(user),
	}, nil
}
