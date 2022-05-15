package handler

import (
	"context"
	"errors"
	"github.com/nuntiodev/block-proto/go_block"
)

func (h *defaultHandler) DeleteText(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error) {
	if req.Text.Id == go_block.LanguageCode_INVALID_LANGUAGE_CODE {
		return nil, errors.New("cannot delete text with invalid language code id")
	}
	config, err := h.repository.Config(ctx, req.Namespace)
	if err != nil {
		return nil, err
	}
	namespaceConfig, err := config.GetNamespaceConfig(ctx)
	if err != nil {
		return nil, err
	}
	if namespaceConfig.DefaultLanguage == req.Text.Id {
		return nil, errors.New("cannot delete default language text")
	}
	text, err := h.repository.Text(ctx, req.Namespace)
	if err != nil {
		return nil, err
	}
	if err := text.Delete(ctx, req.Text.Id); err != nil {
		return nil, err
	}
	return &go_block.UserResponse{}, nil
}
