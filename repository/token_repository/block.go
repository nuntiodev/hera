package token_repository

import (
	"context"
	"errors"
	"github.com/nuntiodev/nuntio-user-block/models"
	"time"

	"github.com/nuntiodev/block-proto/go_block"
	"go.mongodb.org/mongo-driver/bson"
)

func (r *mongodbRepository) Block(ctx context.Context, token *go_block.Token) (*models.Token, error) {
	if token == nil {
		return nil, errors.New("token is nil")
	} else if token.Id == "" {
		return nil, errors.New("missing required token id")
	}
	expiresAt := time.Now().Add(time.Hour * 48)
	mongoUpdate := bson.M{
		"$set": bson.M{
			"blocked":    true,
			"blocked_at": time.Now(),
			"expires_at": expiresAt, // tokens should expire after a day, after being blocked
		},
	}
	result := r.collection.FindOneAndUpdate(
		ctx,
		bson.M{"_id": token.Id},
		mongoUpdate,
	)
	if err := result.Err(); err != nil {
		return nil, err
	}
	var resp models.Token
	if err := result.Decode(&resp); err != nil {
		return nil, err
	}
	if err := r.crypto.Decrypt(&resp); err != nil {
		return nil, err
	}
	// set updated fields
	resp.Blocked = true
	resp.BlockedAt = time.Now()
	resp.ExpiresAt = expiresAt
	return &resp, nil
}
