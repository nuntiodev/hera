package measurement_repository

import (
	"context"
	"github.com/nuntiodev/block-proto/go_block"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"os"
	"time"
)

var (
	activeMeasurementExpiresAt = time.Hour * 24 * 3 // save for 3 days
	enableActiveRepository     = false
)

type defaultMeasurementRepository struct {
	userActiveMeasurementCollection  *mongo.Collection
	userActiveHistoryCollection      *mongo.Collection
	namespaceActiveHistoryCollection *mongo.Collection
	zapLog                           *zap.Logger
}

type MeasurementRepository interface {
	RecordActive(ctx context.Context, measurement *go_block.ActiveMeasurement) (*go_block.ActiveMeasurement, error)
	GetActiveMeasurement(ctx context.Context, measurement *go_block.ActiveMeasurement) (*go_block.ActiveMeasurement, error)
	GetNamespaceActiveHistory(ctx context.Context, year int32) (*go_block.ActiveHistory, error)
	GetUserActiveHistory(ctx context.Context, userId string) (*go_block.ActiveHistory, error)
}

func initialize() error {
	activeMeasurementExpiresAtString, ok := os.LookupEnv("ACTIVE_MEASUREMENT_EXPIRES_AT")
	if ok {
		dur, err := time.ParseDuration(activeMeasurementExpiresAtString)
		if err == nil {
			activeMeasurementExpiresAt = dur
		}
	}
	return nil
}
