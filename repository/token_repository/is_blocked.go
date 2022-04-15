package token_repository

import (
	"context"
	"errors"

	"github.com/io-nuntio/block-proto/go_block"
	"go.mongodb.org/mongo-driver/bson"
)

func (r *mongodbRepository) IsBlocked(ctx context.Context, token *go_block.Token) (bool, error) {
	if token == nil {
		return false, errors.New("token is nil")
	} else if token.Id == "" {
		return false, errors.New("missing required id")
	}
	filter := bson.M{"blocked": true, "_id": token.Id}
	if err := r.collection.FindOne(ctx, filter).Err(); err == nil {
		return true, nil
	}
	return false, nil
}
