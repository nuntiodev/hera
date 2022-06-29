package handler

import (
	"context"

	"github.com/nuntiodev/hera-sdks/go_hera"
)

/*
	Heartbeat - this method checks if the application is live and returns a heartbeat if it is.
*/
func (h *defaultHandler) Heartbeat(ctx context.Context, req *go_hera.HeraRequest) (resp *go_hera.HeraResponse, err error) {
	if err := h.repository.Liveness(ctx); err != nil {
		return nil, err
	}
	return &go_hera.HeraResponse{}, nil
}
