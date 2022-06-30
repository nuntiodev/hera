package handler

import (
	"context"
	"github.com/nuntiodev/hera/repository/config_repository"
	"github.com/nuntiodev/hera/repository/user_repository"
	"golang.org/x/sync/errgroup"

	"github.com/nuntiodev/hera-sdks/go_hera"
)

/*
	ListUsers - this method return a list of users and a count of how many there are total in the database.
*/
func (h *defaultHandler) ListUsers(ctx context.Context, req *go_hera.HeraRequest) (resp *go_hera.HeraResponse, err error) {
	var (
		userRepository   user_repository.UserRepository
		configRepository config_repository.ConfigRepository
		users            []*go_hera.User
		config           *go_hera.Config
		count            int64
		errGroup         = &errgroup.Group{}
	)
	userRepository, err = h.repository.UserRepositoryBuilder().SetNamespace(req.Namespace).Build(ctx)
	if err != nil {
		return nil, err
	}
	// async action 1  - get all users with filter.
	errGroup.Go(func() (err error) {
		users, err = userRepository.List(ctx, req.Query)
		return err
	})
	// async action 2  - get a count of all users in db.
	errGroup.Go(func() (err error) {
		count, err = userRepository.Count(ctx)
		return err
	})
	// async action 3 - get namespace config
	errGroup.Go(func() (err error) {
		configRepository, err = h.repository.ConfigRepositoryBuilder().SetNamespace(req.Namespace).Build(ctx)
		if err != nil {
			return err
		}
		config, err = configRepository.Get(ctx)
		return err
	})
	if err = errGroup.Wait(); err != nil {
		return nil, err
	}
	return &go_hera.HeraResponse{
		Users:  users,
		Amount: count,
		Config: config,
	}, nil
}
