package handler

import (
	"context"
	"time"

	"github.com/nuntiodev/block-proto/go_block"
)

func (h *defaultHandler) UserActiveHistory(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error) {
	measurements, err := h.repository.Measurements(ctx, req.Namespace)
	if err != nil {
		return nil, err
	}
	resp, err := measurements.GetUserActiveHistory(ctx, int32(time.Now().Year()), req.ActiveMeasurement.UserId)
	if err != nil {
		return nil, err
	}
	return &go_block.UserResponse{
		ActiveHistory: resp,
	}, nil
}
