package token_repository

import (
	"context"
	"errors"
	"github.com/nuntiodev/hera/models"
	"strings"
	"time"

	"github.com/nuntiodev/hera-proto/go_hera"
	ts "google.golang.org/protobuf/types/known/timestamppb"
)

func (r *mongodbRepository) Create(ctx context.Context, token *go_hera.Token) error {
	// validate fields
	if token == nil {
		return errors.New("token is nil")
	} else if token.Id == "" {
		return errors.New("missing required token id")
	} else if token.UserId == "" {
		return errors.New("missing required user id")
	} else if token.ExpiresAt == nil || token.ExpiresAt.IsValid() == false {
		return errors.New("missing required token expires at")
	} else if token.ExpiresAt.AsTime().Sub(time.Now()).Seconds() < 0 {
		return errors.New("expires at cannot be in the past")
	} else if token.Type == go_hera.TokenType_TOKEN_TYPE_INVALID {
		return errors.New("invalid token type")
	}
	// prepare fields
	token.Id = strings.TrimSpace(token.Id)
	token.UserId = strings.TrimSpace(token.UserId)
	if token.DeviceInfo == "" {
		token.DeviceInfo = "Unknown"
	}
	token.Blocked = false
	token.CreatedAt = ts.Now()
	token.UsedAt = ts.Now()
	// convert
	create := models.ProtoTokenToToken(token)
	if err := r.crypto.Encrypt(create); err != nil {
		return err
	}
	_, err := r.collection.InsertOne(ctx, create)
	if err != nil {
		return err
	}
	return nil
}
