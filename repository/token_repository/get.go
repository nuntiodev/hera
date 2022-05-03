package token_repository

import (
	"context"
	"errors"
	"github.com/nuntiodev/block-proto/go_block"
	"go.mongodb.org/mongo-driver/bson"
)

func (t *mongodbRepository) Get(ctx context.Context, token *go_block.Token) (*go_block.Token, error) {
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
	if resp.InternalEncryptionLevel > 0 && len(t.internalEncryptionKeys) >= resp.InternalEncryptionLevel {
		if err := t.DecryptToken(&resp); err != nil {
			return nil, err
		}
	}
	return TokenToProtoToken(&resp), nil
}
