package handler

import (
	"context"
	"github.com/nuntiodev/nuntio-user-block/repository/config_repository"
	"github.com/nuntiodev/nuntio-user-block/repository/user_repository"
	"golang.org/x/sync/errgroup"

	"github.com/nuntiodev/block-proto/go_block"
)

/*
	DeleteNamespace - this method deletes an entire namespace. This includes deleting all users and the namespace config.
*/
func (h *defaultHandler) DeleteNamespace(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error) {
	var (
		userRepo   user_repository.UserRepository
		configRepo config_repository.ConfigRepository
		errGroup   = &errgroup.Group{}
		err        error
	)
	// async action 1 - delete all users
	errGroup.Go(func() error {
		userRepo, err = h.repository.Users().SetNamespace(req.Namespace).Build(ctx)
		if err != nil {
			return err
		}
		return userRepo.DeleteAll(ctx)
	})
	// async action 2 - delete config
	errGroup.Go(func() error {
		configRepo, err = h.repository.Config(ctx, req.Namespace)
		if err != nil {
			return err
		}
		return configRepo.Delete(ctx)
	})
	if err = errGroup.Wait(); err != nil {
		return &go_block.UserResponse{}, err
	}
	err = errGroup.Wait()
	return &go_block.UserResponse{}, err
}
