package measurement_repository

import (
	"context"
	"errors"
	"github.com/nuntiodev/block-proto/go_block"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	ts "google.golang.org/protobuf/types/known/timestamppb"
	"k8s.io/utils/pointer"
	"time"
)

func (dmr *defaultMeasurementRepository) RecordActive(ctx context.Context, measurement *go_block.ActiveMeasurement) (*go_block.ActiveMeasurement, error) {
	if measurement == nil {
		return nil, errors.New("")
	} else if measurement.Id == "" {
		return nil, errors.New("missing required id")
	} else if measurement.UserId == "" {
		return nil, errors.New("missing required user id")
	} else if measurement.Seconds == 0 {
		return nil, errors.New("measurement is 0")
	}
	// set default fields
	measurement.CreatedAt = ts.Now()
	measurement.ExpiresAt = ts.New(time.Now().Add(activeMeasurementExpiresAt))
	// create in active measurement collection
	create := ProtoActiveMeasurementToActiveMeasurement(measurement)
	if _, err := dmr.userActiveMeasurementCollection.InsertOne(ctx, create); err != nil {
		return nil, err
	}
	// create in history collections
	// check if it already exists in user history
	now := time.Now()
	year := int32(now.Year())
	month := int32(now.Month())
	userActiveHistory, err := dmr.GetUserActiveHistory(ctx, year, measurement.UserId)
	if err != nil {
		userActiveHistory = &go_block.ActiveHistory{}
		userActiveHistory.UserId = measurement.Id
		userActiveHistory.Year = year
	}
	if _, ok := userActiveHistory.Data[month]; !ok {
		userActiveHistory.Data = map[int32]*go_block.ActiveHistoryData{
			month: {
				Seconds: 0,
				Points:  0,
				From:    map[string]*go_block.CityHistoryMap{},
			},
		}
	}
	userActiveHistory.Data[month].Seconds += measurement.Seconds
	userActiveHistory.Data[month].Points += 1
	userMongoUpdate := bson.M{
		"$set": bson.M{
			"_id":     userActiveHistory.Year,
			"user_id": userActiveHistory.UserId,
			"data":    userActiveHistory.Data,
		},
	}
	if _, err := dmr.userActiveHistoryCollection.UpdateOne(ctx, bson.M{"_id": year}, userMongoUpdate, &options.UpdateOptions{Upsert: pointer.BoolPtr(true)}); err != nil {
		return nil, err
	}
	// now do the same for the namespace collection
	namespaceActiveHistory, err := dmr.GetNamespaceActiveHistory(ctx, year)
	if err != nil {
		namespaceActiveHistory = &go_block.ActiveHistory{}
		namespaceActiveHistory.Year = year
	}
	if _, ok := namespaceActiveHistory.Data[month]; !ok {
		namespaceActiveHistory.Data = map[int32]*go_block.ActiveHistoryData{
			month: {
				Seconds: 0,
				Points:  0,
				From:    map[string]*go_block.CityHistoryMap{},
			},
		}
	}
	namespaceActiveHistory.Data[month].Seconds += measurement.Seconds
	namespaceActiveHistory.Data[month].Points += 1
	if measurement.From != nil && measurement.From.CountryCode != "" {
		if val, ok := namespaceActiveHistory.Data[month].From[measurement.From.CountryCode]; val == nil || !ok {
			namespaceActiveHistory.Data[month].From[measurement.From.CountryCode] = &go_block.CityHistoryMap{
				CityAmount: map[string]int32{},
			}
		}
		namespaceActiveHistory.Data[month].From[measurement.From.CountryCode].CityAmount[measurement.From.City] += 1
	}
	namespaceMongoUpdate := bson.M{
		"$set": bson.M{
			"_id":     namespaceActiveHistory.Year,
			"user_id": namespaceActiveHistory.UserId,
			"data":    namespaceActiveHistory.Data,
		},
	}
	if _, err := dmr.namespaceActiveHistoryCollection.UpdateOne(ctx, bson.M{"_id": year}, namespaceMongoUpdate, &options.UpdateOptions{Upsert: pointer.BoolPtr(true)}); err != nil {
		return nil, err
	}
	return measurement, nil
}
