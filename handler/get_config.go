package handler

import (
	"context"
	"github.com/nuntiodev/block-proto/go_block"
	"github.com/nuntiodev/nuntio-user-block/repository/config_repository"
)

/*
	GetConfig - this method returns a config for a specific namespace.
*/
func (h *defaultHandler) GetConfig(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error) {
	var (
		configRepo config_repository.ConfigRepository
		config     *go_block.Config
		err        error
	)
	configRepo, err = h.repository.Config(ctx, req.Namespace)
	if err != nil {
		return nil, err
	}
	config, err = configRepo.GetNamespaceConfig(ctx)
	return &go_block.UserResponse{
		Config: config,
	}, err
}
