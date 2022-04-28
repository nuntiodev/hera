package measurement_repository

import (
	"context"
	"errors"
	"github.com/nuntiodev/block-proto/go_block"
	"go.mongodb.org/mongo-driver/bson"
	"strings"
)

func (dmr *defaultMeasurementRepository) GetUserActiveHistory(ctx context.Context, year int32, userId string) (*go_block.ActiveHistory, error) {
	if userId == "" {
		return nil, errors.New("missing required user id")
	}
	// prepare
	userId = strings.TrimSpace(userId)
	filter := bson.M{"user_id": getUserHash(userId), "year": year}
	resp := ActiveHistory{}
	if err := dmr.userActiveHistoryCollection.FindOne(ctx, filter).Decode(&resp); err != nil {
		return nil, err
	}
	return ActiveHistoryToProtoActiveHistory(&resp), nil
}
