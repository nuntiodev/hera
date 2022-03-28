package token_repository

import (
	"context"
	"errors"
)

func (r *mongoRepository) BlockToken(ctx context.Context, token *Token) (*Token, error) {
	if token == nil {
		return nil, errors.New("token is nil")
	} else if token.ExpiresAt == 0 {
		return nil, errors.New("token expired at is empty")
	}
	if _, err := r.collection.InsertOne(ctx, token); err != nil {
		return nil, err
	}
	return token, nil
}
