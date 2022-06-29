package token_repository

import (
	"context"
	"errors"
	"github.com/nuntiodev/hera/models"
	"strings"
	"time"

	"github.com/nuntiodev/hera-sdks/go_hera"
	"go.mongodb.org/mongo-driver/bson"
)

func (r *mongodbRepository) UpdateUsedAt(ctx context.Context, token *go_hera.Token) error {
	if token == nil {
		return errors.New("token is nil")
	} else if token.Id == "" {
		return errors.New("missing required token id")
	}
	token.Id = strings.TrimSpace(token.Id)
	token.UserId = strings.TrimSpace(token.UserId)
	if token.DeviceInfo == "" {
		token.DeviceInfo = "Unknown"
	}
	update := models.ProtoTokenToToken(&go_hera.Token{
		DeviceInfo:   token.DeviceInfo,
		LoggedInFrom: token.LoggedInFrom,
	})
	// encrypt data
	if err := r.crypto.Encrypt(update); err != nil {
		return err
	}
	mongoUpdate := bson.M{
		"$set": bson.M{
			"used_at":        time.Now(),
			"device_info":    update.DeviceInfo,
			"logged_in_from": update.LoggedInFrom,
		},
	}
	if _, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": token.Id},
		mongoUpdate,
	); err != nil {
		return err
	}
	return nil
}
