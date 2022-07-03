package handler

import (
	"context"

	"github.com/nuntiodev/hera-sdks/go_hera"
	"github.com/nuntiodev/hera/repository/config_repository"
	"github.com/nuntiodev/hera/repository/user_repository"
)

/*
	CreateUser - this method creates a user in the database with a valid config.
*/
func (h *defaultHandler) CreateUser(ctx context.Context, req *go_hera.HeraRequest) (resp *go_hera.HeraResponse, err error) {
	var (
		configRepository config_repository.ConfigRepository
		userRepository   user_repository.UserRepository
		config           *go_hera.Config
		user             *go_hera.User
	)
	configRepository, err = h.repository.ConfigRepositoryBuilder().SetNamespace(req.Namespace).Build(ctx)
	if err != nil {
		return nil, err
	}
	config, err = configRepository.Get(ctx)
	if err != nil {
		return nil, err
	}
	// build repository used to create user in the database
	h.repository.UserRepositoryBuilder()
	userRepository, err = h.repository.UserRepositoryBuilder().SetNamespace(req.Namespace).WithPasswordValidation(config.ValidatePassword).Build(ctx)
	if err != nil {
		return nil, err
	}
	user, err = userRepository.Create(ctx, req.User)
	if err != nil {
		return nil, err
	}
	return &go_hera.HeraResponse{
		User: user,
	}, nil
}
