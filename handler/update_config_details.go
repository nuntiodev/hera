package handler

import (
	"context"
	"github.com/nuntiodev/block-proto/go_block"
	"github.com/nuntiodev/nuntio-user-block/repository/config_repository"
)

/*
	UpdateConfigDetails - this method updates a namespace's config details such as name and logo.
*/
func (h *defaultHandler) UpdateConfigDetails(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error) {
	var (
		configRepo config_repository.ConfigRepository
		err        error
	)
	configRepo, err = h.repository.Config(ctx, req.Namespace)
	if err != nil {
		return nil, err
	}
	config, err := configRepo.UpdateDetails(ctx, req.Config)
	return &go_block.UserResponse{
		Config: config,
	}, err
}
