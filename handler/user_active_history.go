package handler

import (
	"context"
	"github.com/nuntiodev/nuntio-user-block/repository/measurement_repository"
	"time"

	"github.com/nuntiodev/block-proto/go_block"
)

/*
	UserActiveHistory - this method is used to fetch a user's active history (eg. information about average and total screen time)
*/
func (h *defaultHandler) UserActiveHistory(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error) {
	var (
		measurementRepo measurement_repository.MeasurementRepository
		activeHistory   *go_block.ActiveHistory
		err             error
	)
	measurementRepo, err = h.repository.Measurements(ctx, req.Namespace)
	if err != nil {
		return nil, err
	}
	activeHistory, _, err = measurementRepo.GetUserActiveHistory(ctx, int32(time.Now().Year()), req.ActiveMeasurement.UserId)
	return &go_block.UserResponse{
		ActiveHistory: activeHistory,
	}, err
}
