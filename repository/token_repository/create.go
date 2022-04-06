package token_repository

import (
	"context"
	"errors"
	"strings"
	"time"
)

func (r *mongodbRepository) Create(ctx context.Context, token *Token) (*Token, error) {
	// validate fields
	if token == nil {
		return nil, errors.New("token is nil")
	} else if token.Id == "" {
		return nil, errors.New("missing required token id")
	} else if token.UserId == "" {
		return nil, errors.New("missing required user id")
	} else if token.ExpiresAt.IsZero() {
		return nil, errors.New("missing required token expires at")
	} else if token.ExpiresAt.Sub(time.Now()).Seconds() < 0 {
		return nil, errors.New("expires at cannot be in the past")
	}
	// prepare fields
	token.Id = strings.TrimSpace(token.Id)
	token.UserId = strings.TrimSpace(token.UserId)
	token.Device = strings.TrimSpace(token.Device)
	if token.Device == "" {
		token.Device = "Unknown"
	}
	token.Blocked = false
	token.CreatedAt = time.Now()
	token.UsedAt = time.Now()
	if len(r.internalEncryptionKeys) > 0 {
		if err := r.EncryptToken(actionCreate, token); err != nil {
			return nil, err
		}
		token.Encrypted = true
		token.InternalEncryptionLevel = len(r.internalEncryptionKeys)
	}
	_, err := r.collection.InsertOne(ctx, token)
	if err != nil {
		return nil, err
	}
	return token, nil
}
