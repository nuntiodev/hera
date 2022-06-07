package handler

import (
	"context"
	"github.com/nuntiodev/block-proto/go_block"
	"github.com/nuntiodev/nuntio-user-block/models"
	"github.com/nuntiodev/nuntio-user-block/repository/config_repository"
)

/*
	UpdatePassword - this method updates a users password and validates that password if setting is present in config.
*/
func (h *defaultHandler) UpdatePassword(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error) {
	var (
		configRepo config_repository.ConfigRepository
		config     *models.Config
		user       *models.User
		err        error
	)
	// get config
	configRepo, err = h.repository.Config(ctx, req.Namespace, req.EncryptionKey)
	if err != nil {
		return &go_block.UserResponse{}, err
	}
	config, err = configRepo.GetNamespaceConfig(ctx)
	if err != nil {
		return &go_block.UserResponse{}, err
	}
	users, err := h.repository.UserRepositoryBuilder().SetNamespace(req.Namespace).WithPasswordValidation(config.ValidatePassword).Build(ctx)
	if err != nil {
		return &go_block.UserResponse{}, err
	}
	user, err = users.UpdatePassword(ctx, req.User, req.Update)
	if err != nil {
		return &go_block.UserResponse{}, err
	}
	return &go_block.UserResponse{
		User: models.UserToProtoUser(user),
	}, nil
}
