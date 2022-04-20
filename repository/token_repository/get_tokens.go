package token_repository

import (
	"context"
	"errors"
	"github.com/nuntiodev/block-proto/go_block"
	"go.mongodb.org/mongo-driver/bson"
)

func (t *mongodbRepository) GetTokens(ctx context.Context, token *go_block.Token) ([]*go_block.Token, error) {
	if token == nil {
		return nil, errors.New("token is nil")
	} else if token.UserId == "" {
		return nil, errors.New("missing required user id")
	}
	filter := bson.M{"user_id": token.UserId}
	cursor, err := t.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var resp []*go_block.Token
	for cursor.Next(ctx) {
		tempToken := Token{}
		if err := cursor.Decode(&tempToken); err != nil {
			return nil, err
		}
		if tempToken.Encrypted {
			if err := t.DecryptToken(&tempToken); err != nil {
				return nil, err
			}
		}
		resp = append(resp, TokenToProtoToken(&tempToken))
	}
	return resp, nil
}

// Currently, tokens are not upgraded to higher encryption levels when we add another encryption key. this is primarily due to the fact that
// tokens should be short-lived, meaning they will regularly be deleted from the database and is thus not considered a security issue.
// Also, tokens do not contain any harmful information other than name of device which is encrypted.
