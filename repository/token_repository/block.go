package token_repository

import (
	"context"
	"errors"
	"time"

	"github.com/io-nuntio/block-proto/go_block"
	"go.mongodb.org/mongo-driver/bson"
	ts "google.golang.org/protobuf/types/known/timestamppb"
)

func (r *mongodbRepository) Block(ctx context.Context, token *go_block.Token) (*go_block.Token, error) {
	if token == nil {
		return nil, errors.New("token is nil")
	} else if token.Id == "" {
		return nil, errors.New("missing required token id")
	}
	mongoUpdate := bson.M{
		"$set": bson.M{
			"blocked":    true,
			"blocked_at": time.Now(),
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
	token.Blocked = true
	token.BlockedAt = ts.Now()
	return token, nil
}
