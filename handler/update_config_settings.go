package handler

import (
	"context"
	"github.com/nuntiodev/block-proto/go_block"
	"github.com/nuntiodev/nuntio-user-block/models"
	"github.com/nuntiodev/nuntio-user-block/repository/config_repository"
)

/*
	UpdateConfigSettings - this method updates a namespace's config settings such as password and email requirements.
*/
func (h *defaultHandler) UpdateConfigSettings(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error) {
	var (
		configRepo config_repository.ConfigRepository
		err        error
	)
	configRepo, err = h.repository.Config(ctx, req.Namespace, req.EncryptionKey)
	if err != nil {
		return nil, err
	}
	config, err := configRepo.UpdateSettings(ctx, req.Config)
	return &go_block.UserResponse{
		Config: models.ConfigToProtoConfig(config),
	}, err
}
