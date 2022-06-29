package handler

import (
	"context"
	"github.com/nuntiodev/hera-sdks/go_hera"
	"github.com/nuntiodev/hera/repository/config_repository"
)

/*
	UpdateConfig - this method updates a namespace's config such as name and logo, validate password and etc.
*/
func (h *defaultHandler) UpdateConfig(ctx context.Context, req *go_hera.HeraRequest) (resp *go_hera.HeraResponse, err error) {
	var (
		configRepository config_repository.ConfigRepository
	)
	configRepository, err = h.repository.ConfigRepositoryBuilder().SetNamespace(req.Namespace).Build(ctx)
	if err != nil {
		return nil, err
	}
	if err := configRepository.Update(ctx, req.Config); err != nil {
		return nil, err
	}
	return &go_hera.HeraResponse{}, nil
}
