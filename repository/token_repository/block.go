package token_repository

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

func (r *mongodbRepository) Block(ctx context.Context, token *Token) (*Token, error) {
	if token == nil {
		return nil, errors.New("token is nil")
	} else if token.Id == "" {
		return nil, errors.New("missing required token id")
	}
	token.Blocked = true
	token.BlockedAt = time.Now()
	mongoUpdate := bson.M{
		"$set": bson.M{
			"blocked":    token.Blocked,
			"blocked_at": token.BlockedAt,
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
	return token, nil
}
