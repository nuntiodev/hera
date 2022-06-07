package email_repository

import (
	"context"
	"errors"
	"github.com/nuntiodev/block-proto/go_block"
	"github.com/nuntiodev/nuntio-user-block/models"
	"go.mongodb.org/mongo-driver/bson"
)

func (e *defaultEmailRepository) GetAll(ctx context.Context, email *go_block.Email) ([]*models.Email, error) {
	if email == nil {
		return nil, errors.New("email is nil")
	} else if email.Id == "" {
		return nil, errors.New("missing required id")
	}
	prepare(actionGet, email)
	var resp []*models.Email
	cursor, err := e.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	for cursor.Next(ctx) {
		temp := models.Email{}
		if err := cursor.Decode(&temp); err != nil {
			return nil, err
		}
		if err := e.crypto.Decrypt(&temp); err != nil {
			return nil, err
		}
		// check if we should upgrade the encryption level
		if upgradable, _ := e.crypto.Upgradeble(&temp); upgradable {
			if err := e.upgradeEncryptionLevel(ctx, &temp); err != nil {
				return nil, err
			}
		}
		resp = append(resp, &temp)
	}
	return resp, nil
}
