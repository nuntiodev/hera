package token_repository

import (
	"context"
	"errors"
	"github.com/nuntiodev/nuntio-user-block/models"
	"github.com/nuntiodev/x/cryptox"
	"strings"
	"time"

	"github.com/nuntiodev/block-proto/go_block"
	"go.mongodb.org/mongo-driver/bson"
)

func (r *mongodbRepository) UpdateUsedAt(ctx context.Context, token *go_block.Token) (*models.Token, error) {
	if token == nil {
		return nil, errors.New("token is nil")
	} else if token.Id == "" {
		return nil, errors.New("missing required token id")
	}
	token.Id = strings.TrimSpace(token.Id)
	token.UserId = strings.TrimSpace(token.UserId)
	if token.DeviceInfo == "" {
		token.DeviceInfo = "Unknown"
	}
	update := models.ProtoTokenToToken(&go_block.Token{
		DeviceInfo:   token.DeviceInfo,
		LoggedInFrom: token.LoggedInFrom,
	})
	// encrypt data
	if err := r.crypto.Encrypt(update); err != nil {
		return nil, err
	}
	mongoUpdate := bson.M{
		"$set": bson.M{
			"used_at":        time.Now(),
			"device_info":    update.DeviceInfo,
			"logged_in_from": update.LoggedInFrom,
		},
	}
	result := r.collection.FindOneAndUpdate(
		ctx,
		bson.M{"_id": token.Id},
		mongoUpdate,
	)
	if err := result.Err(); err != nil {
		return nil, err
	}
	var resp models.Token
	if err := result.Decode(&resp); err != nil {
		return nil, err
	}
	if err := r.crypto.Decrypt(&resp); err != nil {
		return nil, err
	}
	// set updated fields
	resp.UsedAt = time.Now()
	resp.DeviceInfo = cryptox.Stringx{
		Body:                    token.DeviceInfo,
		InternalEncryptionLevel: update.DeviceInfo.InternalEncryptionLevel,
		ExternalEncryptionLevel: update.DeviceInfo.ExternalEncryptionLevel,
	}
	if token.LoggedInFrom != nil {
		resp.LoggedInFrom = &models.Location{
			Country: cryptox.Stringx{
				Body:                    token.LoggedInFrom.Country,
				InternalEncryptionLevel: update.LoggedInFrom.Country.InternalEncryptionLevel,
				ExternalEncryptionLevel: update.LoggedInFrom.Country.ExternalEncryptionLevel,
			},
			CountryCode: cryptox.Stringx{
				Body:                    token.LoggedInFrom.CountryCode,
				InternalEncryptionLevel: update.LoggedInFrom.CountryCode.InternalEncryptionLevel,
				ExternalEncryptionLevel: update.LoggedInFrom.CountryCode.ExternalEncryptionLevel,
			},
			City: cryptox.Stringx{
				Body:                    token.LoggedInFrom.City,
				InternalEncryptionLevel: update.LoggedInFrom.City.InternalEncryptionLevel,
				ExternalEncryptionLevel: update.LoggedInFrom.City.ExternalEncryptionLevel,
			},
		}
	}
	return &resp, nil
}
