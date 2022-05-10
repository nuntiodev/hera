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
	// get previous token
	get, err := r.Get(ctx, token)
	if err != nil {
		return nil, err
	}
	token.DeviceInfo = strings.TrimSpace(token.DeviceInfo)
	token.Id = strings.TrimSpace(token.Id)
	token.UserId = strings.TrimSpace(token.UserId)
	if token.DeviceInfo == "" {
		token.DeviceInfo = "Unknown"
	}
	update := ProtoTokenToToken(token)
	// encrypt data
	if get.InternalEncryptionLevel > 0 {
		update.InternalEncryptionLevel = int(get.InternalEncryptionLevel)
		if err := r.EncryptToken(actionUpdate, update); err != nil {
			return nil, err
		}
	}
	mongoUpdate := bson.M{
		"$set": bson.M{
			"used_at":        time.Now(),
			"device_info":    update.DeviceInfo,
			"logged_in_from": update.LoggedInFrom,
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
	get.UsedAt = ts.Now()
	get.DeviceInfo = token.DeviceInfo
	get.LoggedInFrom = token.LoggedInFrom
	return get, nil
}
