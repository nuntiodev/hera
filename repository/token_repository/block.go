package token_repository

import (
	"context"
	"errors"
	"time"

	"github.com/nuntiodev/hera-proto/go_hera"
	"go.mongodb.org/mongo-driver/bson"
)

func (r *mongodbRepository) Block(ctx context.Context, token *go_hera.Token) error {
	if token == nil {
		return errors.New("token is nil")
	} else if token.Id == "" {
		return errors.New("missing required token id")
	}
	expiresAt := time.Now().Add(time.Hour * 12)
	mongoUpdate := bson.M{
		"$set": bson.M{
			"blocked":    true,
			"blocked_at": time.Now(),
			"expires_at": expiresAt, // tokens should expire after a day, after being blocked
		},
	}
	if _, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": token.Id},
		mongoUpdate,
	); err != nil {
		return err
	}
	return nil
}
