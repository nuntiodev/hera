package models

import (
	"github.com/nuntiodev/block-proto/go_block"
	"github.com/nuntiodev/x/cryptox"
	ts "google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type Location struct {
	Country     cryptox.Stringx `bson:"country" json:"country"`
	CountryCode cryptox.Stringx `bson:"country_code" json:"country_code"`
	City        cryptox.Stringx `bson:"city" json:"city"`
}

type Token struct {
	Id           string             `bson:"_id" json:"_id"`
	UserId       string             `bson:"user_id" json:"user_id"`
	Blocked      bool               `bson:"blocked" json:"blocked"`
	DeviceInfo   cryptox.Stringx    `bson:"device_info" json:"device_info"`
	LoggedInFrom *Location          `bson:"logged_in_from" json:"logged_in_from"`
	Type         go_block.TokenType `bson:"token_type" json:"token_type"`
	BlockedAt    time.Time          `bson:"blocked_at" json:"blocked_at"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
	UsedAt       time.Time          `bson:"used_at" json:"used_at"`
	ExpiresAt    time.Time          `bson:"expires_at" json:"expires_at"` // unix time
}

func LocationToProto(location *Location) *go_block.Location {
	if location == nil {
		return nil
	}
	return &go_block.Location{
		City:        location.City.Body,
		Country:     location.Country.Body,
		CountryCode: location.CountryCode.Body,
	}
}

func ProtoToLocation(location *go_block.Location) *Location {
	if location == nil {
		return nil
	}
	return &Location{
		City:        cryptox.Stringx{Body: location.City},
		Country:     cryptox.Stringx{Body: location.Country},
		CountryCode: cryptox.Stringx{Body: location.CountryCode},
	}
}

func TokensToProto(tokens []*Token) []*go_block.Token {
	var resp []*go_block.Token
	for _, token := range tokens {
		resp = append(resp, TokenToProtoToken(token))
	}
	return resp
}

func TokenToProtoToken(token *Token) *go_block.Token {
	if token == nil {
		return nil
	}
	location := &go_block.Location{}
	if token.LoggedInFrom != nil {
		location.Country = token.LoggedInFrom.Country.Body
		location.CountryCode = token.LoggedInFrom.CountryCode.Body
		location.City = token.LoggedInFrom.City.Body
	}
	return &go_block.Token{
		Id:           token.Id,
		UserId:       token.UserId,
		Blocked:      token.Blocked,
		DeviceInfo:   token.DeviceInfo.Body,
		LoggedInFrom: location,
		Type:         token.Type,
		BlockedAt:    ts.New(token.BlockedAt),
		CreatedAt:    ts.New(token.CreatedAt),
		UsedAt:       ts.New(token.UsedAt),
		ExpiresAt:    ts.New(token.ExpiresAt),
	}
}

func ProtoTokenToToken(token *go_block.Token) *Token {
	if token == nil {
		return nil
	}
	location := Location{}
	if token.LoggedInFrom != nil {
		location.Country = cryptox.Stringx{Body: token.LoggedInFrom.Country}
		location.CountryCode = cryptox.Stringx{Body: token.LoggedInFrom.CountryCode}
		location.City = cryptox.Stringx{Body: token.LoggedInFrom.City}
	}
	return &Token{
		Id:           token.Id,
		UserId:       token.UserId,
		Blocked:      token.Blocked,
		DeviceInfo:   cryptox.Stringx{Body: token.DeviceInfo},
		LoggedInFrom: &location,
		Type:         token.Type,
		BlockedAt:    token.BlockedAt.AsTime(),
		CreatedAt:    token.CreatedAt.AsTime(),
		UsedAt:       token.UsedAt.AsTime(),
		ExpiresAt:    token.ExpiresAt.AsTime(),
	}
}
