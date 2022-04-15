package handler

import (
	"context"

	"github.com/io-nuntio/block-proto/go_block"
)

func (h *defaultHandler) DeleteBatch(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error) {
	users, err := h.repository.Users().SetNamespace(req.Namespace).Build(ctx)
	if err != nil {
		return &go_block.UserResponse{}, err
	}
	return &go_block.UserResponse{}, users.DeleteBatch(ctx, req.UserBatch)
}
