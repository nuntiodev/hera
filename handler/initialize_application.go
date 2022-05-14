package handler

import (
	"context"
	"github.com/nuntiodev/block-proto/go_block"
)

func (h *defaultHandler) InitializeApplication(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error) {
	config, err := h.repository.Config(ctx, req.Namespace)
	if err != nil {
		return nil, err
	}
	namespaceConfig, err := config.GetNamespaceConfig(ctx)
	if err != nil {
		return nil, err
	}
	text, err := h.repository.Text(ctx, req.Namespace)
	if err != nil {
		return nil, err
	}
	defaultLanguage, err := text.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return &go_block.UserResponse{
		Config: namespaceConfig,
		Texts:  defaultLanguage,
	}, nil
}
