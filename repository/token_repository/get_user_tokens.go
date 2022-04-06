package token_repository

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
)

func (r *mongodbRepository) GetUserTokens(ctx context.Context, token *Token) ([]*Token, error) {
	if token == nil {
		return nil, errors.New("token is nil")
	} else if token.UserId == "" {
		return nil, errors.New("missing required user id")
	}
	filter := bson.M{"user_id": token.UserId}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var resp []*Token
	for cursor.Next(ctx) {
		token := Token{}
		if err := cursor.Decode(&token); err != nil {
			return nil, err
		}
		resp = append(resp, &token)
	}
	return resp, nil
}

// Currently, tokens are not upgraded to higher encryption levels when we add another encryption key. this is primarily due to the fact that
// tokens should be short-lived, meaning they will regularly be deleted from the database and is thus not considered a security issue.
// Also, tokens do not contain any harmful information other than name of device which is encrypted.
