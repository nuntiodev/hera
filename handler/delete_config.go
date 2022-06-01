package handler

import (
	"context"
	"github.com/nuntiodev/block-proto/go_block"
	"github.com/nuntiodev/nuntio-user-block/repository/config_repository"
)

/*
	DeleteConfig - this method deletes a namespace config.
*/
func (h *defaultHandler) DeleteConfig(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error) {
	var (
		configRepo config_repository.ConfigRepository
		err        error
	)
	configRepo, err = h.repository.Config(ctx, req.Namespace)
	if err != nil {
		return nil, err
	}
	return &go_block.UserResponse{}, configRepo.Delete(ctx)
}
