package token_repository

import (
	"context"
	"errors"
	"strings"
)

func (r *mongoRepository) BlockToken(ctx context.Context, token *Token) error {
	if token == nil {
		return errors.New("token is nil")
	} else if token.ExpiresAt == 0 {
		return errors.New("token expired at is empty")
	}
	token.RefreshTokenId = strings.TrimSpace(token.RefreshTokenId)
	token.AccessTokenId = strings.TrimSpace(token.AccessTokenId)
	if _, err := r.collection.InsertOne(ctx, token); err != nil {
		return err
	}
	return nil
}
