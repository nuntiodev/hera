package handler

import (
	"context"
	"errors"
	"github.com/nuntiodev/block-proto/go_block"
	"github.com/nuntiodev/nuntio-user-block/models"
	"github.com/nuntiodev/nuntio-user-block/repository/config_repository"
	"github.com/nuntiodev/nuntio-user-block/repository/user_repository"
)

/*
	Create - this method creates a user in the database with a valid config.
*/
func (h *defaultHandler) Create(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error) {
	var (
		configRepo config_repository.ConfigRepository
		userRepo   user_repository.UserRepository
		config     *models.Config
		user       *models.User
		err        error
	)
	configRepo, err = h.repository.Config(ctx, req.Namespace, req.EncryptionKey)
	if err != nil {
		return &go_block.UserResponse{}, err
	}
	config, err = configRepo.GetNamespaceConfig(ctx)
	if err != nil {
		return &go_block.UserResponse{}, err
	}
	// validate that the action is possible with project config
	// we cannot send an email if the email provider is not enabled
	if !h.emailEnabled && config.RequireEmailVerification {
		return &go_block.UserResponse{}, errors.New("email provider is not enabled and verification email cannot be sent")
	} else if config.RequireEmailVerification && (req.User.Email == "") {
		return &go_block.UserResponse{}, errors.New("require email verification is enabled and user email is empty")
	}
	userRepo, err = h.repository.UserRepositoryBuilder().SetNamespace(req.Namespace).SetEncryptionKey(req.EncryptionKey).WithPasswordValidation(config.ValidatePassword).Build(ctx)
	if err != nil {
		return &go_block.UserResponse{}, err
	}
	user, err = userRepo.Create(ctx, req.User)
	return &go_block.UserResponse{
		User: models.UserToProtoUser(user),
	}, err
}
