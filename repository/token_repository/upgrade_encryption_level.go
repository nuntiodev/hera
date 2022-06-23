package token_repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/nuntiodev/hera/models"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

func (e *mongodbRepository) upgradeEncryptionLevel(ctx context.Context, token *models.Token) error {
	if token == nil {
		return errors.New("token is nil")
	}
	if upgradable, err := e.crypto.Upgradeble(token); err != nil && !upgradable {
		return fmt.Errorf("could not upgrade with err: %v", err)
	}
	copy := *token
	if err := e.crypto.Encrypt(&copy); err != nil {
		return err
	}
	mongoUpdate := bson.M{
		"$set": bson.M{
			"device_info":    copy.DeviceInfo,
			"logged_in_from": copy.LoggedInFrom,
			"updated_at":     time.Now(),
			"encrypted_at":   time.Now(),
		},
	}
	if _, err := e.collection.UpdateOne(
		ctx,
		bson.M{"_id": copy.Id},
		mongoUpdate,
	); err != nil {
		return err
	}
	// update levels
	token.DeviceInfo.InternalEncryptionLevel = copy.DeviceInfo.InternalEncryptionLevel
	token.LoggedInFrom.InternalEncryptionLevel = copy.LoggedInFrom.InternalEncryptionLevel
	return nil
}
