package handler

import (
	"context"
	"github.com/nuntiodev/block-proto/go_block"
	"github.com/nuntiodev/nuntio-user-block/models"
	"github.com/nuntiodev/nuntio-user-block/repository/config_repository"
)

/*
	UpdateConfig - this method updates a namespace's config such as name and logo, validate password and etc.
*/
func (h *defaultHandler) UpdateConfig(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error) {
	var (
		configRepo config_repository.ConfigRepository
		err        error
	)
	configRepo, err = h.repository.Config(ctx, req.Namespace, req.EncryptionKey)
	if err != nil {
		return nil, err
	}
	config, err := configRepo.Update(ctx, req.Config)
	return &go_block.UserResponse{
		Config: models.ConfigToProtoConfig(config),
	}, err
}
