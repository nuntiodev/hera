package handler

import (
	"context"
	"errors"
	"fmt"
	"github.com/nuntiodev/nuntio-user-block/models"
	"github.com/nuntiodev/nuntio-user-block/repository/config_repository"
	"github.com/nuntiodev/nuntio-user-block/repository/user_repository"

	"github.com/nuntiodev/block-proto/go_block"
)

/*
	UpdateEmail - this method updates a users email.
*/
func (h *defaultHandler) UpdateEmail(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error) {
	var (
		configRepo config_repository.ConfigRepository
		userRepo   user_repository.UserRepository
		config     *models.Config
		err        error
	)
	configRepo, err = h.repository.Config(ctx, req.Namespace, req.EncryptionKey)
	if err != nil {
		return &go_block.UserResponse{}, fmt.Errorf("could not build config repository with err: %v", err)
	}
	config, err = configRepo.GetNamespaceConfig(ctx)
	if err != nil {
		return &go_block.UserResponse{}, fmt.Errorf("could not get namespace config with err: %v", err)
	}
	// we cannot send an email if the email provider is not enabled
	if !h.emailEnabled && config.RequireEmailVerification {
		return &go_block.UserResponse{}, errors.New("email provider is not enabled and verification email cannot be sent")
	} else if config.RequireEmailVerification && req.Update.Email == "" {
		return &go_block.UserResponse{}, errors.New("require email is enabled and email is empty")
	}
	userRepo, err = h.repository.UserRepositoryBuilder().SetNamespace(req.Namespace).SetEncryptionKey(req.EncryptionKey).Build(ctx)
	if err != nil {
		return &go_block.UserResponse{}, fmt.Errorf("could not build user repository with err: %v", err)
	}
	updatedUser, err := userRepo.UpdateEmail(ctx, req.User, req.Update)
	if err != nil {
		return &go_block.UserResponse{}, fmt.Errorf("could not update user email with err: %v", err)
	}
	return &go_block.UserResponse{
		User: models.UserToProtoUser(updatedUser),
	}, nil
}
