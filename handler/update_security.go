package handler

import (
	"context"
	"github.com/nuntiodev/nuntio-user-block/models"
	"github.com/nuntiodev/nuntio-user-block/repository/user_repository"

	"github.com/nuntiodev/block-proto/go_block"
)

/*
	UpdateSecurity - this method updates a user's security settings. If encrypted the user will be decrypted.
	If not, the user will be encrypted.
	todo: move this setting to config.
*/
func (h *defaultHandler) UpdateSecurity(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error) {
	var (
		userRepo user_repository.UserRepository
		user     *models.User
		err      error
	)
	userRepo, err = h.repository.UserRepositoryBuilder().SetNamespace(req.Namespace).SetEncryptionKey(req.EncryptionKey).Build(ctx)
	if err != nil {
		return &go_block.UserResponse{}, err
	}
	user, err = userRepo.UpdateSecurity(ctx, req.User)
	return &go_block.UserResponse{
		User: models.UserToProtoUser(user),
	}, err
}
