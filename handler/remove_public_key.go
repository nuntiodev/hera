package handler

import (
	"context"
	"github.com/nuntiodev/hera-sdks/go_hera"
	"github.com/nuntiodev/hera/repository/config_repository"
)

func (h *defaultHandler) RemovePublicKey(ctx context.Context, req *go_hera.HeraRequest) (resp *go_hera.HeraResponse, err error) {
	var (
		configRepository config_repository.ConfigRepository
	)
	return &go_hera.HeraResponse{}, nil // not implemented yet
	configRepository, err = h.repository.ConfigRepositoryBuilder().SetNamespace(req.Namespace).Build(ctx)
	if err != nil {
		return nil, err
	}
	// todo: decrypt everything using the private key before removing it
	if err = configRepository.RemovePublicKey(ctx); err != nil {
		return nil, err
	}
	return nil, nil
}
