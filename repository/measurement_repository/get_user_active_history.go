package measurement_repository

import (
	"context"
	"errors"
	"github.com/nuntiodev/block-proto/go_block"
	"go.mongodb.org/mongo-driver/bson"
	"strings"
)

func (dmr *defaultMeasurementRepository) GetUserActiveHistory(ctx context.Context, year int32, userId string) (*go_block.ActiveHistory, bool, error) {
	if userId == "" {
		return nil, false, errors.New("missing required user id")
	}
	// prepare
	userId = strings.TrimSpace(userId)
	filter := bson.M{"user_id": getUserHash(userId), "year": year}
	resp := ActiveHistory{}
	res := dmr.userActiveHistoryCollection.FindOne(ctx, filter)
	if err := res.Err(); err != nil {
		return nil, false, err
	}
	if err := res.Decode(&resp); err != nil {
		return nil, true, err
	}
	return ActiveHistoryToProtoActiveHistory(&resp), true, nil
}
