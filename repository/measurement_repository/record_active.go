package measurement_repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/nuntiodev/block-proto/go_block"
	"go.mongodb.org/mongo-driver/bson"
	ts "google.golang.org/protobuf/types/known/timestamppb"
	"strings"
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
	// prepare
	measurement.Id = strings.TrimSpace(measurement.Id)
	measurement.UserId = strings.TrimSpace(measurement.UserId)
	// set default fields
	measurement.CreatedAt = ts.Now()
	measurement.ExpiresAt = ts.New(time.Now().Add(activeMeasurementExpiresAt))
	// create in active measurement collection
	create := ProtoActiveMeasurementToActiveMeasurement(measurement)
	if _, err := dmr.userActiveMeasurementCollection.InsertOne(ctx, create); err != nil {
		return nil, err
	}
	// create in user history collection
	now := time.Now()
	year := int32(now.Year())
	month := int32(now.Month())
	userActiveHistory, alreadyCreated, err := dmr.GetUserActiveHistory(ctx, year, measurement.UserId)
	if alreadyCreated && err != nil {
		return nil, fmt.Errorf("could not decode user active history with err: %v", err)
	}
	if err != nil {
		userActiveHistory = &go_block.ActiveHistory{}
		userActiveHistory.Year = year
		// set user id hash instead of just id; this is more secure
		userActiveHistory.UserId = getUserHash(measurement.UserId)
	}
	if _, ok := userActiveHistory.Data[month]; !ok {
		userActiveHistory.Data[month] = &go_block.ActiveHistoryData{
			Seconds: 0,
			Points:  0,
			From:    map[string]*go_block.CityHistoryMap{},
			Device:  map[string]int32{},
		}
	}
	// make sure data is initialized
	if userActiveHistory.Data[month].From == nil {
		userActiveHistory.Data[month].From = map[string]*go_block.CityHistoryMap{}
	}
	if userActiveHistory.Data[month].Device == nil {
		userActiveHistory.Data[month].Device = map[string]int32{}
	}
	// set optional data
	if measurement.From != nil && measurement.From.CountryCode != "" {
		if val, ok := userActiveHistory.Data[month].From[measurement.From.CountryCode]; val == nil || !ok {
			// country does not exist in map yet; initialize it to 0
			userActiveHistory.Data[month].From[measurement.From.CountryCode] = &go_block.CityHistoryMap{
				CityAmount: map[string]int32{
					measurement.From.City: 0,
				},
			}
		}
		userActiveHistory.Data[month].From[measurement.From.CountryCode].CityAmount[measurement.From.City] += 1
	}
	if measurement.Device != go_block.Platform_INVALID_PLATFORM {
		userActiveHistory.Data[month].Device[measurement.Device.String()] += 1
	}
	// set required data
	userActiveHistory.Data[month].Seconds += measurement.Seconds
	userActiveHistory.Data[month].Points += 1
	userMongoUpdate := bson.M{
		"$set": bson.M{
			"data": userActiveHistory.Data,
		},
	}
	if alreadyCreated {
		if _, err := dmr.userActiveHistoryCollection.UpdateOne(ctx, bson.M{"user_id": userActiveHistory.UserId}, userMongoUpdate); err != nil {
			return nil, err
		}
	} else {
		update := ProtoActiveHistoryToActiveHistory(userActiveHistory)
		if _, err := dmr.userActiveHistoryCollection.InsertOne(ctx, update); err != nil {
			return nil, err
		}
	}
	// now do the same for the namespace collection
	return measurement, nil
}
