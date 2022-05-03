package token_repository

import (
	"github.com/nuntiodev/block-proto/go_block"
	ts "google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type Location struct {
	Country     string `bson:"country" json:"country"`
	CountryCode string `bson:"country_code" json:"country_code"`
	City        string `bson:"city" json:"city"`
}

type Token struct {
	Id                      string             `bson:"_id" json:"_id"`
	UserId                  string             `bson:"user_id" json:"user_id"`
	Blocked                 bool               `bson:"blocked" json:"blocked"`
	DeviceInfo              string             `bson:"device_info" json:"device_info"`
	LoggedInFrom            Location           `bson:"logged_in_from" json:"logged_in_from"`
	Type                    go_block.TokenType `bson:"token_type" json:"token_type"`
	BlockedAt               time.Time          `bson:"blocked_at" json:"blocked_at"`
	CreatedAt               time.Time          `bson:"created_at" json:"created_at"`
	UsedAt                  time.Time          `bson:"used_at" json:"used_at"`
	ExpiresAt               time.Time          `bson:"expires_at" json:"expires_at"` // unix time
	InternalEncryptionLevel int                `bson:"internal_encryption_level" json:"internal_encryption_level"`
}

func TokenToProtoToken(token *Token) *go_block.Token {
	if token == nil {
		return nil
	}
	location := &go_block.Location{
		Country:     token.LoggedInFrom.Country,
		CountryCode: token.LoggedInFrom.CountryCode,
		City:        token.LoggedInFrom.City,
	}
	return &go_block.Token{
		Id:                      token.Id,
		UserId:                  token.UserId,
		Blocked:                 token.Blocked,
		DeviceInfo:              token.DeviceInfo,
		LoggedInFrom:            location,
		Type:                    token.Type,
		BlockedAt:               ts.New(token.BlockedAt),
		CreatedAt:               ts.New(token.CreatedAt),
		UsedAt:                  ts.New(token.UsedAt),
		ExpiresAt:               ts.New(token.ExpiresAt),
		InternalEncryptionLevel: int32(token.InternalEncryptionLevel),
	}
}

func ProtoTokenToToken(token *go_block.Token) *Token {
	if token == nil {
		return nil
	}
	location := Location{}
	if token.LoggedInFrom != nil {
		location.Country = token.LoggedInFrom.Country
		location.CountryCode = token.LoggedInFrom.CountryCode
		location.City = token.LoggedInFrom.City
	}
	return &Token{
		Id:                      token.Id,
		UserId:                  token.UserId,
		Blocked:                 token.Blocked,
		DeviceInfo:              token.DeviceInfo,
		LoggedInFrom:            location,
		Type:                    token.Type,
		BlockedAt:               token.BlockedAt.AsTime(),
		CreatedAt:               token.CreatedAt.AsTime(),
		UsedAt:                  token.UsedAt.AsTime(),
		ExpiresAt:               token.ExpiresAt.AsTime(),
		InternalEncryptionLevel: int(token.InternalEncryptionLevel),
	}
}
