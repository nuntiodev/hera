package measurement_repository

import (
	"context"
	"errors"
	"fmt"
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
		fmt.Println("get one 1")
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
	userUpdateOrCreate := ProtoActiveHistoryToActiveHistory(userActiveHistory)
	if _, err := dmr.userActiveHistoryCollection.UpdateOne(ctx, bson.M{"_id": year}, userUpdateOrCreate, &options.UpdateOptions{Upsert: pointer.BoolPtr(true)}); err != nil {
		fmt.Println("get one 2")
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
		namespaceActiveHistory.Data[month].From[measurement.From.CountryCode].CityAmount[measurement.From.City] += 1
	}
	namespaceUpdateOrCreate := ProtoActiveHistoryToActiveHistory(userActiveHistory)
	if _, err := dmr.namespaceActiveHistoryCollection.UpdateOne(ctx, bson.M{"_id": year}, namespaceUpdateOrCreate, &options.UpdateOptions{Upsert: pointer.BoolPtr(true)}); err != nil {
		fmt.Println("get one 3")
		return nil, err
	}
	return measurement, nil
}
