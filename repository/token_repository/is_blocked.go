package token_repository

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
)

func (r *mongoRepository) IsBlocked(ctx context.Context, token *Token) error {
	if token == nil {
		return errors.New("token is nil")
	}
	filter := bson.M{"_id": token.Id}
	if err := r.collection.FindOne(ctx, filter).Err(); err == nil {
		return errors.New("found blocked token")
	}
	return nil
}
