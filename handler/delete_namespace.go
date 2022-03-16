package handler

import (
	"context"
	"github.com/softcorp-io/block-proto/go_block"
)

func (h *defaultHandler) DeleteNamespace(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error) {
	return &go_block.UserResponse{}, h.repository.UserRepository.DeleteNamespace(ctx, req.Namespace)
}
