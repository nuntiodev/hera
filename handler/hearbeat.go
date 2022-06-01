package handler

import (
	"context"

	"github.com/nuntiodev/block-proto/go_block"
)

/*
	Heartbeat - this method checks if the application is live and returns a heartbeat if it is.
*/
func (h *defaultHandler) Heartbeat(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error) {
	return &go_block.UserResponse{}, h.repository.Liveness(ctx)
}
