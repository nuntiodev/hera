package handler

import (
	"context"
	"github.com/softcorp-io/block-proto/go_block"
)

func (h *defaultHandler) DeleteBatch(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error) {
	return &go_block.UserResponse{}, h.repository.UserRepository.DeleteBatch(ctx, req.UserBatch, req.Namespace)
}
