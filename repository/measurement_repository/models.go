package measurement_repository

import (
	"github.com/nuntiodev/block-proto/go_block"
	ts "google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

// ActiveMeasurement struct for how to measure user time
type ActiveMeasurement struct {
	Id        string    `bson:"_id" json:"_id"`
	UserId    string    `bson:"user_id" json:"user_id"`
	Seconds   int32     `bson:"seconds" json:"seconds"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	ExpiresAt time.Time `bson:"expires_at" json:"expires_at"`
}

type ActiveHistoryData struct {
	Seconds int32            `bson:"seconds" json:"seconds"`
	Points  int32            `bson:"points" json:"points"`
	From    map[string]int32 `bson:"from" json:"from"`
}

// ActiveHistory keeps map over data over a year and maps data from month to ActiveHistoryData
type ActiveHistory struct {
	Year   int32                       `bson:"_id" json:"_id"`
	UserId string                      `bson:"user_id" json:"user_id"`
	Data   map[int32]ActiveHistoryData `bson:"data" json:"data"`
}

func ActiveMeasurementToProtoActiveMeasurement(active *ActiveMeasurement) *go_block.ActiveMeasurement {
	if active == nil {
		return nil
	}
	return &go_block.ActiveMeasurement{
		Id:        active.Id,
		UserId:    active.UserId,
		Seconds:   active.Seconds,
		CreatedAt: ts.New(active.CreatedAt),
	}
}

func ProtoActiveMeasurementToActiveMeasurement(active *go_block.ActiveMeasurement) *ActiveMeasurement {
	if active == nil {
		return nil
	}
	return &ActiveMeasurement{
		Id:        active.Id,
		UserId:    active.UserId,
		Seconds:   active.Seconds,
		CreatedAt: active.CreatedAt.AsTime(),
	}
}

func ActiveHistoryToProtoActiveHistory(history *ActiveHistory) *go_block.ActiveHistory {
	if history == nil {
		return nil
	}
	data := map[int32]*go_block.ActiveHistoryData{}
	for k, v := range history.Data {
		data[k] = &go_block.ActiveHistoryData{
			Seconds: v.Seconds,
			Points:  v.Points,
			From:    v.From,
		}
	}
	return &go_block.ActiveHistory{
		Year:   history.Year,
		UserId: history.UserId,
		Data:   data,
	}
}

func ProtoActiveHistoryToActiveHistory(history *go_block.ActiveHistory) *ActiveHistory {
	if history == nil {
		return nil
	}
	data := map[int32]ActiveHistoryData{}
	for k, v := range history.Data {
		data[k] = ActiveHistoryData{
			Seconds: v.Seconds,
			Points:  v.Points,
			From:    v.From,
		}
	}
	return &ActiveHistory{
		Year:   history.Year,
		UserId: history.UserId,
		Data:   data,
	}
}
