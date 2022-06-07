package email_repository

import (
	"context"
	"errors"
	"github.com/nuntiodev/block-proto/go_block"
	"github.com/nuntiodev/nuntio-user-block/models"
	"go.mongodb.org/mongo-driver/bson"
)

func (e *defaultEmailRepository) Get(ctx context.Context, email *go_block.Email) (*models.Email, error) {
	if email == nil {
		return nil, errors.New("email is nil")
	} else if email.Id == "" {
		return nil, errors.New("missing required id")
	}
	prepare(actionGet, email)
	resp := models.Email{}
	result := e.collection.FindOne(ctx, bson.M{"_id": email.Id})
	if err := result.Err(); err != nil {
		return nil, err
	}
	if err := result.Decode(&resp); err != nil {
		return nil, err
	}
	if err := e.crypto.Decrypt(&resp); err != nil {
		return nil, err
	}
	// check if we should upgrade the encryption level
	if upgradable, _ := e.crypto.Upgradeble(&resp); upgradable {
		if err := e.upgradeEncryptionLevel(ctx, &resp); err != nil {
			return nil, err
		}
	}
	return &resp, nil
}
