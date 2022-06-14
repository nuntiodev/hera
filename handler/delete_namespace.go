package handler

import (
	"context"
	"github.com/nuntiodev/block-proto/go_block"
)

/*
	DeleteNamespace - this method deletes an entire namespace. This includes deleting all users and the namespace config.
*/
func (h *defaultHandler) DeleteNamespace(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error) {
	return &go_block.UserResponse{}, h.repository.DropDatabase(ctx, req.Namespace)
}
