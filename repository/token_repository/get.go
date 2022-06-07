package token_repository

import (
	"context"
	"errors"
	"github.com/nuntiodev/block-proto/go_block"
	"github.com/nuntiodev/nuntio-user-block/models"
	"go.mongodb.org/mongo-driver/bson"
)

func (t *mongodbRepository) Get(ctx context.Context, token *go_block.Token) (*models.Token, error) {
	if token == nil {
		return nil, errors.New("token is nil")
	} else if token.Id == "" {
		return nil, errors.New("missing required user id")
	}
	filter := bson.M{"_id": token.Id}
	resp := models.Token{}
	if err := t.collection.FindOne(ctx, filter).Decode(&resp); err != nil {
		return nil, err
	}
	if err := t.crypto.Decrypt(&resp); err != nil {
		return nil, err
	}
	// check if we should upgrade the encryption level
	if upgradable, _ := t.crypto.Upgradeble(&resp); upgradable {
		if err := t.upgradeEncryptionLevel(ctx, &resp); err != nil {
			return nil, err
		}
	}
	return &resp, nil
}
