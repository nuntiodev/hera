package token_repository

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/nuntiodev/block-proto/go_block"
	"go.mongodb.org/mongo-driver/bson"
	ts "google.golang.org/protobuf/types/known/timestamppb"
)

func (r *mongodbRepository) UpdateUsedAt(ctx context.Context, token *go_block.Token) (*go_block.Token, error) {
	if token == nil {
		return nil, errors.New("token is nil")
	} else if token.Id == "" {
		return nil, errors.New("missing required token id")
	}
	token.DeviceInfo = strings.TrimSpace(token.DeviceInfo)
	token.Id = strings.TrimSpace(token.Id)
	token.UserId = strings.TrimSpace(token.UserId)
	if token.DeviceInfo == "" {
		token.DeviceInfo = "Unknown"
	}
	if token.LoggedInFrom == "" {
		token.LoggedInFrom = "Unknown"
	}
	mongoUpdate := bson.M{
		"$set": bson.M{
			"used_at":        time.Now(),
			"device":         token.DeviceInfo,
			"logged_in_from": token.LoggedInFrom,
		},
	}
	updateResult, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": token.Id},
		mongoUpdate,
	)
	if err != nil {
		return nil, err
	}
	if updateResult.MatchedCount == 0 {
		return nil, errors.New("could not find token")
	}
	// set updated fields
	token.UsedAt = ts.Now()
	return token, nil
}
