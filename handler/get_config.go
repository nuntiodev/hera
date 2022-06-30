package handler

import (
	"context"
	"github.com/nuntiodev/hera-sdks/go_hera"
	"github.com/nuntiodev/hera/repository/config_repository"
)

/*
	GetConfig - this method returns a config for a specific namespace.
*/
func (h *defaultHandler) GetConfig(ctx context.Context, req *go_hera.HeraRequest) (resp *go_hera.HeraResponse, err error) {
	var (
		configRepository config_repository.ConfigRepository
		config           *go_hera.Config
	)
	configRepository, err = h.repository.ConfigRepositoryBuilder().SetNamespace(req.Namespace).Build(ctx)
	if err != nil {
		return nil, err
	}
	config, err = configRepository.Get(ctx)
	if err != nil {
		return nil, err
	}
	return &go_hera.HeraResponse{
		Config: config,
	}, nil
}
