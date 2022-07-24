package handler

import (
	"context"
	"github.com/nuntiodev/hera-sdks/go_hera"
	"github.com/nuntiodev/hera/hash"
	"github.com/nuntiodev/hera/repository/config_repository"
	"github.com/nuntiodev/hera/repository/user_repository"
)

/*
	UpdateUserPassword - this method updates a users password and validates that password if setting is present in config.
*/
func (h *defaultHandler) UpdateUserPassword(ctx context.Context, req *go_hera.HeraRequest) (resp *go_hera.HeraResponse, err error) {
	var (
		configRepository config_repository.ConfigRepository
		userRepository   user_repository.UserRepository
		config           *go_hera.Config
		user             *go_hera.User
	)
	// get config
	configRepository, err = h.repository.ConfigRepositoryBuilder().SetNamespace(req.Namespace).Build(ctx)
	if err != nil {
		return nil, err
	}
	config, err = configRepository.Get(ctx)
	if err != nil {
		return nil, err
	}
	userRepository, err = h.repository.UserRepositoryBuilder().SetNamespace(req.Namespace).SetHasher(hash.New(config)).WithPasswordValidation(config.ValidatePassword).Build(ctx)
	if err != nil {
		return nil, err
	}
	if err = userRepository.UpdatePassword(ctx, req.User, req.UserUpdate); err != nil {
		return nil, err
	}
	return &go_hera.HeraResponse{
		User: user,
	}, nil
}
