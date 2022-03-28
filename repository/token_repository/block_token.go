package token_repository

import (
	"context"
	"errors"
)

func (r *mongoRepository) BlockToken(ctx context.Context, token *Token) error {
	if token == nil {
		return errors.New("token is nil")
	} else if token.ExpiresAt == 0 {
		return errors.New("token expired at is empty")
	}
	if _, err := r.collection.InsertOne(ctx, token); err != nil {
		return err
	}
	return nil
}
