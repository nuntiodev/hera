package handler

import (
	"context"
	"github.com/nuntiodev/nuntio-user-block/models"
	"github.com/nuntiodev/nuntio-user-block/repository/config_repository"
	"github.com/nuntiodev/nuntio-user-block/repository/user_repository"
	"golang.org/x/sync/errgroup"

	"github.com/nuntiodev/block-proto/go_block"
)

/*
	GetAll - this method return all users and a count of how many there are.
*/
func (h *defaultHandler) GetAll(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error) {
	var (
		userRepo   user_repository.UserRepository
		configRepo config_repository.ConfigRepository
		users      []*models.User
		config     *models.Config
		count      int64
		errGroup   = &errgroup.Group{}
		err        error
	)
	userRepo, err = h.repository.UserRepositoryBuilder().SetNamespace(req.Namespace).SetEncryptionKey(req.EncryptionKey).Build(ctx)
	if err != nil {
		return &go_block.UserResponse{}, err
	}
	// async action 1  - get all users with filter.
	errGroup.Go(func() error {
		users, err = userRepo.GetAll(ctx, req.Filter)
		return err
	})
	// async action 2  - get a count of all users in db.
	errGroup.Go(func() error {
		count, err = userRepo.Count(ctx)
		return err
	})
	// async action 3 - get namespace config
	errGroup.Go(func() error {
		configRepo, err = h.repository.Config(ctx, req.Namespace, req.EncryptionKey)
		if err != nil {
			return err
		}
		config, err = configRepo.GetNamespaceConfig(ctx)
		return err
	})
	err = errGroup.Wait()
	return &go_block.UserResponse{
		Users:       models.UsersToProto(users),
		UsersAmount: count,
		Config:      models.ConfigToProtoConfig(config),
	}, err
}
