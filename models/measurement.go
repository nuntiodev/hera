package models

import (
	"time"

	"github.com/nuntiodev/block-proto/go_block"
	ts "google.golang.org/protobuf/types/known/timestamppb"
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
	Seconds int32                               `bson:"seconds" json:"seconds"`
	Points  int32                               `bson:"points" json:"points"`
	From    map[string]*go_block.CityHistoryMap `bson:"from" json:"from"`     // from country to city to int
	Dau     map[int32]string                    `bson:"dau" json:"dau"`       // day of month int to user hash
	Device  map[go_block.Platform]int32         `bson:"device" json:"device"` //
}

// ActiveHistory keeps map over data over a year and maps data from month to ActiveHistoryData
type ActiveHistory struct {
	Year   int32                        `bson:"year" json:"year"`
	UserId string                       `bson:"user_id" json:"user_id"`
	Data   map[int32]*ActiveHistoryData `bson:"data" json:"data"`
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
			Dau:     v.Dau, // Dau maps day of month to hash of user id -> dau is amount of keys in map
			Device:  deviceToProtoDevice(v.Device),
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
	data := map[int32]*ActiveHistoryData{}
	for k, v := range history.Data {
		data[k] = &ActiveHistoryData{
			Seconds: v.Seconds,
			Points:  v.Points,
			From:    v.From,
			Dau:     v.Dau,
			Device:  protoDeviceToDevice(v.Device),
		}
	}
	return &ActiveHistory{
		Year:   history.Year,
		UserId: history.UserId,
		Data:   data,
	}
}

func protoDeviceToDevice(device map[string]int32) map[go_block.Platform]int32 {
	if device == nil {
		return map[go_block.Platform]int32{}
	}
	resp := map[go_block.Platform]int32{}
	for k, v := range device {
		switch k {
		case go_block.Platform_IOS.String():
			resp[go_block.Platform_IOS] = v
		case go_block.Platform_ANDROID.String():
			resp[go_block.Platform_ANDROID] = v
		case go_block.Platform_WEB.String():
			resp[go_block.Platform_WEB] = v
		case go_block.Platform_MACOS.String():
			resp[go_block.Platform_MACOS] = v
		case go_block.Platform_LINUX.String():
			resp[go_block.Platform_LINUX] = v
		case go_block.Platform_WINDOWS.String():
			resp[go_block.Platform_WINDOWS] = v
		default:
			resp[go_block.Platform_INVALID_PLATFORM] = v
		}
	}
	return resp
}

func deviceToProtoDevice(device map[go_block.Platform]int32) map[string]int32 {
	if device == nil {
		return map[string]int32{}
	}
	resp := map[string]int32{}
	for k, v := range device {
		switch k {
		case go_block.Platform_IOS:
			resp[go_block.Platform_IOS.String()] = v
		case go_block.Platform_ANDROID:
			resp[go_block.Platform_ANDROID.String()] = v
		case go_block.Platform_WEB:
			resp[go_block.Platform_WEB.String()] = v
		case go_block.Platform_MACOS:
			resp[go_block.Platform_MACOS.String()] = v
		case go_block.Platform_LINUX:
			resp[go_block.Platform_LINUX.String()] = v
		case go_block.Platform_WINDOWS:
			resp[go_block.Platform_WINDOWS.String()] = v
		default:
			resp[go_block.Platform_INVALID_PLATFORM.String()] = v
		}
	}
	return resp
}
