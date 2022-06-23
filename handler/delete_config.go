package handler

import (
	"context"
	"github.com/nuntiodev/hera-proto/go_hera"
	"github.com/nuntiodev/hera/repository/config_repository"
)

/*
	DeleteConfig - this method deletes a namespace config.
*/
func (h *defaultHandler) DeleteConfig(ctx context.Context, req *go_hera.HeraRequest) (resp *go_hera.HeraResponse, err error) {
	var (
		configRepository config_repository.ConfigRepository
	)
	configRepository, err = h.repository.ConfigRepositoryBuilder().SetNamespace(req.Namespace).Build(ctx)
	if err != nil {
		return nil, err
	}
	return nil, configRepository.Delete(ctx)
}
