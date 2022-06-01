package handler

import (
	"context"
	"github.com/nuntiodev/nuntio-user-block/repository/measurement_repository"
	"time"

	"github.com/nuntiodev/block-proto/go_block"
)

/*
	NamespaceActiveHistory - this method returns statistics about average screen time / total screen time from users provided by the client SDKs.
*/
func (h *defaultHandler) NamespaceActiveHistory(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error) {
	var (
		measurementRepo measurement_repository.MeasurementRepository
		activeHistory   *go_block.ActiveHistory
		err             error
	)
	measurementRepo, err = h.repository.Measurements(ctx, req.Namespace)
	if err != nil {
		return nil, err
	}
	activeHistory, err = measurementRepo.GetNamespaceActiveHistory(ctx, int32(time.Now().Year()))
	if err != nil {
		return nil, err
	}
	return &go_block.UserResponse{
		ActiveHistory: activeHistory,
	}, nil
}
