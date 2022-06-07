package handler

import (
	"context"
	"github.com/nuntiodev/nuntio-user-block/models"
	"github.com/nuntiodev/nuntio-user-block/repository/measurement_repository"

	"github.com/nuntiodev/block-proto/go_block"
)

/*
	RecordActiveMeasurement - record active measurement is used by the client SDKs to record average and total screen time for every user.
*/
func (h *defaultHandler) RecordActiveMeasurement(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error) {
	var (
		measurementRepo   measurement_repository.MeasurementRepository
		activeMeasurement *models.ActiveMeasurement
		err               error
	)
	measurementRepo, err = h.repository.Measurements(ctx, req.Namespace)
	if err != nil {
		return nil, err
	}
	activeMeasurement, err = measurementRepo.RecordActive(ctx, req.ActiveMeasurement)
	return &go_block.UserResponse{
		ActiveMeasurement: models.ActiveMeasurementToProtoActiveMeasurement(activeMeasurement),
	}, err
}
