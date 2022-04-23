package handler

import (
	"context"

	"github.com/nuntiodev/block-proto/go_block"
)

func (h *defaultHandler) RecordActiveMeasurement(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error) {
	measurements, err := h.repository.Measurements(ctx, req.Namespace)
	if err != nil {
		return nil, err
	}
	resp, err := measurements.RecordActive(ctx, req.ActiveMeasurement)
	if err != nil {
		return nil, err
	}
	return &go_block.UserResponse{
		ActiveMeasurement: resp,
	}, nil
}
