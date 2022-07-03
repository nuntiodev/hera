package handler

import (
	"context"
	"fmt"
	"github.com/nuntiodev/hera-sdks/go_hera"
	"github.com/nuntiodev/hera/repository/config_repository"
	"github.com/nuntiodev/hera/repository/token_repository"
	"github.com/nuntiodev/hera/repository/user_repository"
	"github.com/nuntiodev/x/pointerx"
	"golang.org/x/sync/errgroup"
)

/*
	CreateNamespace - this method creates a sets up a namespace for a new client.
	This includes creating a database, collections and indexes.
*/
func (h *defaultHandler) CreateNamespace(ctx context.Context, req *go_hera.HeraRequest) (resp *go_hera.HeraResponse, err error) {
	var (
		userRepository   user_repository.UserRepository
		configRepository config_repository.ConfigRepository
		tokenRepository  token_repository.TokenRepository
		errGroup         = &errgroup.Group{}
	)
	// async action 1 - setup user repository and create test user
	errGroup.Go(func() error {
		var err error
		userRepository, err = h.repository.UserRepositoryBuilder().SetNamespace(req.Namespace).Build(ctx)
		if err != nil {
			return fmt.Errorf("could not build user repository with err: %v", err)
		}
		if err := userRepository.BuildIndexes(ctx); err != nil {
			return err
		}
		// create test user
		if _, err = userRepository.Create(ctx, &go_hera.User{
			FirstName: pointerx.StringPtr("Test"),
			LastName:  pointerx.StringPtr("User"),
			Email:     pointerx.StringPtr("test@user.io"),
		}); err != nil {
			return err
		}
		return nil
	})
	// async action 2 - setup config repository and create default config
	errGroup.Go(func() error {
		var err error
		// create initial config
		configRepository, err = h.repository.ConfigRepositoryBuilder().SetNamespace(req.Namespace).Build(ctx)
		if err != nil {
			return fmt.Errorf("could not build config with err: %v", err)
		}
		if err = configRepository.Create(ctx, req.Config); err != nil {
			return err
		}
		return nil
	})
	// async action 3 - setup tokens repository
	errGroup.Go(func() error {
		var err error
		tokenRepository, err = h.repository.TokenRepositoryBuilder().SetNamespace(req.Namespace).Build(ctx)
		if err != nil {
			return err
		}
		if err := tokenRepository.BuildIndexes(ctx); err != nil {
			return err
		}
		return nil
	})
	if err = errGroup.Wait(); err != nil {
		return nil, err
	}
	return &go_hera.HeraResponse{}, nil
}
