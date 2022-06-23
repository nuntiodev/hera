package token_repository

import (
	"context"
	"errors"
	"github.com/nuntiodev/hera-proto/go_hera"
	"github.com/nuntiodev/hera/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (t *mongodbRepository) GetTokens(ctx context.Context, token *go_hera.Token) ([]*models.Token, error) {
	if token == nil {
		return nil, errors.New("token is nil")
	} else if token.UserId == "" {
		return nil, errors.New("missing required user id")
	}
	sortOptions := (&options.FindOptions{}).SetSort(bson.D{{"used_at", -1}, {"_id", -1}})
	filter := bson.M{"user_id": token.UserId}
	cursor, err := t.collection.Find(ctx, filter, sortOptions)
	if err != nil {
		return nil, err
	}
	var resp []*models.Token
	for cursor.Next(ctx) {
		tempToken := models.Token{}
		if err := cursor.Decode(&tempToken); err != nil {
			return nil, err
		}
		if err := t.crypto.Decrypt(&tempToken); err != nil {
			return nil, err
		}
		// check if we should upgrade the encryption level
		if upgradable, _ := t.crypto.Upgradeble(&tempToken); upgradable {
			if err := t.upgradeEncryptionLevel(ctx, &tempToken); err != nil {
				return nil, err
			}
		}
		resp = append(resp, &tempToken)
	}
	return resp, nil
}
