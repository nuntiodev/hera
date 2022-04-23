package measurement_repository

import (
	"context"
	"errors"
	"github.com/nuntiodev/block-proto/go_block"
	"go.mongodb.org/mongo-driver/bson"
)

func (dmr *defaultMeasurementRepository) GetNamespaceActiveHistory(ctx context.Context, year int32) (*go_block.ActiveHistory, error) {
	if year == 0 {
		return nil, errors.New("missing required year")
	}
	filter := bson.M{"_id": year}
	resp := ActiveHistory{}
	if err := dmr.namespaceActiveHistoryCollection.FindOne(ctx, filter).Decode(&resp); err != nil {
		return nil, err
	}
	return ActiveHistoryToProtoActiveHistory(&resp), nil
}
