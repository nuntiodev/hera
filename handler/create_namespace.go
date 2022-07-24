package handler

import (
	"context"
	"fmt"
	"github.com/nuntiodev/hera/hash"
	"github.com/nuntiodev/x/pointerx"

	"github.com/nuntiodev/hera-sdks/go_hera"
	"github.com/nuntiodev/hera/repository/config_repository"
	"github.com/nuntiodev/hera/repository/token_repository"
	"github.com/nuntiodev/hera/repository/user_repository"
	"golang.org/x/sync/errgroup"
)

/*
	CreateNamespace creates and sets up a namespace for a new client.
	This includes creating a database, collections, indexes and etc.
*/
func (h *defaultHandler) CreateNamespace(ctx context.Context, req *go_hera.HeraRequest) (resp *go_hera.HeraResponse, err error) {
	var (
		userRepository   user_repository.UserRepository
		configRepository config_repository.ConfigRepository
		tokenRepository  token_repository.TokenRepository
		config           *go_hera.Config
		errGroup         = &errgroup.Group{}
	)
	// async action 1 - setup user & config repository, create default config and create test user
	errGroup.Go(func() (err error) {
		// create initial config
		configRepository, err = h.repository.ConfigRepositoryBuilder().SetNamespace(req.Namespace).Build(ctx)
		if err != nil {
			return fmt.Errorf("could not build config with err: %v", err)
		}
		config, err = configRepository.Create(ctx, req.Config)
		if err != nil {
			return err
		}
		userRepository, err = h.repository.UserRepositoryBuilder().SetNamespace(req.Namespace).SetHasher(hash.New(config)).Build(ctx)
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
			Password:  &go_hera.Hash{Body: "MySecretPassword1234!*"},
		}); err != nil {
			return err
		}
		return nil
	})
	// async action 2 - setup tokens repository
	errGroup.Go(func() (err error) {
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
	return &go_hera.HeraResponse{
		Config: config,
	}, nil
}
