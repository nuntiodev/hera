package measurement_repository

import (
	"context"
	"errors"
	"github.com/nuntiodev/block-proto/go_block"
	"github.com/nuntiodev/nuntio-user-block/models"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/sync/errgroup"
)

func (dmr *defaultMeasurementRepository) GetNamespaceActiveHistory(ctx context.Context, year int32) (*models.ActiveHistory, error) {
	if year == 0 {
		return nil, errors.New("missing required year")
	}
	resp := models.ActiveHistory{
		Year: year,
		Data: map[int32]*models.ActiveHistoryData{},
	}
	cursor, err := dmr.userActiveHistoryCollection.Find(ctx, bson.M{"year": year})
	if err != nil {
		return nil, err
	}
	g := new(errgroup.Group)
	for cursor.Next(ctx) {
		temp := models.ActiveHistory{}
		if err := cursor.Decode(&temp); err != nil {
			return nil, err
		}
		g.Go(func() error {
			for month, data := range temp.Data {
				resp.Data[month].Seconds += data.Seconds
				resp.Data[month].Points += data.Points
				// validate resp device map is initialized
				if resp.Data[month].Device == nil {
					resp.Data[month].Device = map[go_block.Platform]int32{}
				}
				for device, sum := range temp.Data[month].Device {
					resp.Data[month].Device[device] += sum
				}
				// validate resp dau map is initialized
				if resp.Data[month].Dau == nil {
					resp.Data[month].Dau = map[int32]string{}
				}
				for day, dau := range temp.Data[month].Dau {
					resp.Data[month].Dau[day] += dau
				}
				// validate resp device map is initialized
				if resp.Data[month].From == nil {
					resp.Data[month].From = map[string]*go_block.CityHistoryMap{}
				}
				for country, cityMap := range temp.Data[month].From {
					for city, amount := range cityMap.CityAmount {
						// validate resp city map is initialized
						if resp.Data[month].From == nil {
							resp.Data[month].From[country].CityAmount = map[string]int32{}
						}
						resp.Data[month].From[country].CityAmount[city] += amount
					}
				}
			}
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return nil, err
	}
	resp.Year = year
	return &resp, nil
}
