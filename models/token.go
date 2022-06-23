package models

import (
	"github.com/nuntiodev/hera-proto/go_hera"
	"github.com/nuntiodev/x/cryptox"
	ts "google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type Token struct {
	Id           string            `bson:"_id" json:"_id"`
	UserId       string            `bson:"user_id" json:"user_id"`
	Blocked      bool              `bson:"blocked" json:"blocked"`
	DeviceInfo   cryptox.Stringx   `bson:"device_info" json:"device_info"`
	LoggedInFrom cryptox.Stringx   `bson:"logged_in_from" json:"logged_in_from"`
	Type         go_hera.TokenType `bson:"token_type" json:"token_type"`
	BlockedAt    time.Time         `bson:"blocked_at" json:"blocked_at"`
	CreatedAt    time.Time         `bson:"created_at" json:"created_at"`
	UsedAt       time.Time         `bson:"used_at" json:"used_at"`
	ExpiresAt    time.Time         `bson:"expires_at" json:"expires_at"` // unix time
}

func TokensToProto(tokens []*Token) []*go_hera.Token {
	var resp []*go_hera.Token
	for _, token := range tokens {
		resp = append(resp, TokenToProtoToken(token))
	}
	return resp
}

func TokenToProtoToken(token *Token) *go_hera.Token {
	if token == nil {
		return nil
	}
	return &go_hera.Token{
		Id:           token.Id,
		UserId:       token.UserId,
		Blocked:      token.Blocked,
		DeviceInfo:   token.DeviceInfo.Body,
		LoggedInFrom: token.LoggedInFrom.Body,
		Type:         token.Type,
		BlockedAt:    ts.New(token.BlockedAt),
		CreatedAt:    ts.New(token.CreatedAt),
		UsedAt:       ts.New(token.UsedAt),
		ExpiresAt:    ts.New(token.ExpiresAt),
	}
}

func ProtoTokenToToken(token *go_hera.Token) *Token {
	if token == nil {
		return nil
	}
	return &Token{
		Id:           token.Id,
		UserId:       token.UserId,
		Blocked:      token.Blocked,
		DeviceInfo:   cryptox.Stringx{Body: token.DeviceInfo},
		LoggedInFrom: cryptox.Stringx{Body: token.LoggedInFrom},
		Type:         token.Type,
		BlockedAt:    token.BlockedAt.AsTime(),
		CreatedAt:    token.CreatedAt.AsTime(),
		UsedAt:       token.UsedAt.AsTime(),
		ExpiresAt:    token.ExpiresAt.AsTime(),
	}
}
