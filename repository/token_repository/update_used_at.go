package token_repository

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

func (r *mongodbRepository) UpdateUsedAt(ctx context.Context, token *Token) (*Token, error) {
	if token == nil {
		return nil, errors.New("token is nil")
	} else if token.Id == "" {
		return nil, errors.New("missing required token id")
	}
	token.UsedAt = time.Now()
	mongoUpdate := bson.M{
		"$set": bson.M{
			"usedAt": token.UsedAt,
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
