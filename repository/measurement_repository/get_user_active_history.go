package measurement_repository

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/nuntiodev/block-proto/go_block"
	"go.mongodb.org/mongo-driver/bson"
)

func (dmr *defaultMeasurementRepository) GetUserActiveHistory(ctx context.Context, year int32, userId string) (*go_block.ActiveHistory, error) {
	if userId == "" {
		return nil, errors.New("missing required user id")
	}
	fmt.Println(userId)
	hash := sha256.New()
	hash.Write([]byte(userId))
	userShaHash := string(hash.Sum(nil))
	fmt.Println(userShaHash)
	filter := bson.M{"user_id": userShaHash, "year": year}
	resp := ActiveHistory{}
	if err := dmr.userActiveHistoryCollection.FindOne(ctx, filter).Decode(&resp); err != nil {
		return nil, err
	}
	return ActiveHistoryToProtoActiveHistory(&resp), nil
}
