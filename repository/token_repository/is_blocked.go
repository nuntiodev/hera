package token_repository

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strings"
)

func (r *mongoRepository) IsBlocked(ctx context.Context, token *Token) error {
	if token == nil {
		return errors.New("token is nil")
	}
	token.RefreshTokenId = strings.TrimSpace(token.RefreshTokenId)
	token.AccessTokenId = strings.TrimSpace(token.AccessTokenId)
	var filter interface{}
	if token.AccessTokenId != "" && token.RefreshTokenId != "" {
		filter = bson.D{
			{"$or", bson.A{
				bson.D{{"access_token_id", primitive.Regex{Pattern: token.AccessTokenId, Options: ""}}},
				bson.D{{"refresh_token_id", primitive.Regex{Pattern: token.RefreshTokenId, Options: ""}}},
			},
			},
		}
	} else if token.AccessTokenId != "" {
		filter = bson.M{"access_token_id": token.AccessTokenId}
	} else if token.RefreshTokenId != "" {
		filter = bson.M{"refresh_token_id": token.RefreshTokenId}
	}
	if err := r.collection.FindOne(ctx, filter).Err(); err == nil {
		return errors.New("found blocked token")
	}
	return nil
}
