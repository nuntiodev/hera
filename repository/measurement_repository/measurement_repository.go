package measurement_repository

import (
	"context"
	"github.com/nuntiodev/block-proto/go_block"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"time"
)

const (
	expiresAfterIndex = "expires-after-measurement-index"
)

var (
	activeMeasurementExpiresAt = time.Hour * 24 * 3 // save for 3 days
	enableActiveRepository     = false
)

type defaultMeasurementRepository struct {
	userActiveMeasurementCollection *mongo.Collection
	userActiveHistoryCollection     *mongo.Collection
}

type MeasurementRepository interface {
	RecordActive(ctx context.Context, measurement *go_block.ActiveMeasurement) (*go_block.ActiveMeasurement, error)
	GetNamespaceActiveHistory(ctx context.Context, year int32) (*go_block.ActiveHistory, error)
	GetUserActiveHistory(ctx context.Context, year int32, userId string) (*go_block.ActiveHistory, bool, error)
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

func newMongodbMeasurementRepository(ctx context.Context, userActiveMeasurementCollection, userActiveHistoryCollection *mongo.Collection) (*defaultMeasurementRepository, error) {
	expiresAtIndexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "expires_at", Value: 1},
		},
		Options: options.Index().SetExpireAfterSeconds(0).SetName(expiresAfterIndex),
	}
	if _, err := userActiveMeasurementCollection.Indexes().CreateOne(ctx, expiresAtIndexModel); err != nil {
		return nil, err
	}
	userIdIndexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "user_id", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	}
	if _, err := userActiveHistoryCollection.Indexes().CreateOne(ctx, userIdIndexModel); err != nil {
		return nil, err
	}
	return &defaultMeasurementRepository{
		userActiveMeasurementCollection: userActiveMeasurementCollection,
		userActiveHistoryCollection:     userActiveHistoryCollection,
	}, nil
}

func New(ctx context.Context, userActiveMeasurementCollection, userActiveHistoryCollection *mongo.Collection) (MeasurementRepository, error) {
	return newMongodbMeasurementRepository(ctx, userActiveMeasurementCollection, userActiveHistoryCollection)
}
