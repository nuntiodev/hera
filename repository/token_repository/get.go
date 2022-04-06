package token_repository

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
)

func (t *mongodbRepository) Get(ctx context.Context, token *Token) (*Token, error) {
	if token == nil {
		return nil, errors.New("token is nil")
	} else if token.Id == "" {
		return nil, errors.New("missing required user id")
	}
	filter := bson.M{"_id": token.Id}
	resp := Token{}
	if err := t.collection.FindOne(ctx, filter).Decode(&resp); err != nil {
		return nil, err
	}
	if resp.Encrypted {
		if err := t.DecryptToken(&resp); err != nil {
			return nil, err
		}
	}
	return &resp, nil
}
