package token_repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/nuntiodev/nuntio-user-block/models"
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
	if copy.LoggedInFrom == nil {
		copy.LoggedInFrom = &models.Location{}
	}
	mongoUpdate := bson.M{
		"$set": bson.M{
			"device_info":                 copy.DeviceInfo,
			"logged_in_from.country":      copy.LoggedInFrom.Country,
			"logged_in_from.country_code": copy.LoggedInFrom.CountryCode,
			"logged_in_from.city":         copy.LoggedInFrom.City,
			"updated_at":                  time.Now(),
			"encrypted_at":                time.Now(),
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
	if token.LoggedInFrom == nil {
		token.LoggedInFrom = &models.Location{}
	}
	token.DeviceInfo.InternalEncryptionLevel = copy.DeviceInfo.InternalEncryptionLevel
	token.DeviceInfo.ExternalEncryptionLevel = copy.DeviceInfo.ExternalEncryptionLevel
	token.LoggedInFrom.Country.InternalEncryptionLevel = copy.LoggedInFrom.Country.InternalEncryptionLevel
	token.LoggedInFrom.Country.ExternalEncryptionLevel = copy.LoggedInFrom.Country.ExternalEncryptionLevel
	token.LoggedInFrom.CountryCode.InternalEncryptionLevel = copy.LoggedInFrom.CountryCode.InternalEncryptionLevel
	token.LoggedInFrom.CountryCode.ExternalEncryptionLevel = copy.LoggedInFrom.CountryCode.ExternalEncryptionLevel
	token.LoggedInFrom.City.InternalEncryptionLevel = copy.LoggedInFrom.City.InternalEncryptionLevel
	token.LoggedInFrom.City.ExternalEncryptionLevel = copy.LoggedInFrom.City.ExternalEncryptionLevel
	return nil
}
