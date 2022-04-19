package token_repository

import (
	"github.com/nuntiodev/block-proto/go_block"
	ts "google.golang.org/protobuf/types/known/timestamppb"
)

func TokenToProtoToken(token *Token) *go_block.Token {
	if token == nil {
		return nil
	}
	return &go_block.Token{
		Id:         token.Id,
		UserId:     token.UserId,
		Blocked:    token.Blocked,
		DeviceInfo: token.Device,
		Location: &go_block.Location{
			Country: token.Location,
		},
		BlockedAt:               ts.New(token.BlockedAt),
		CreatedAt:               ts.New(token.CreatedAt),
		UsedAt:                  ts.New(token.UsedAt),
		ExpiresAt:               ts.New(token.ExpiresAt),
		Encrypted:               token.Encrypted,
		InternalEncryptionLevel: int32(token.InternalEncryptionLevel),
	}
}

func ProtoTokenToToken(token *go_block.Token) *Token {
	if token == nil {
		return nil
	}
	location := ""
	if token.Location != nil {
		location = token.Location.Country
	}
	return &Token{
		Id:                      token.Id,
		UserId:                  token.UserId,
		Blocked:                 token.Blocked,
		Device:                  token.DeviceInfo,
		Location:                location,
		BlockedAt:               token.BlockedAt.AsTime(),
		CreatedAt:               token.CreatedAt.AsTime(),
		UsedAt:                  token.UsedAt.AsTime(),
		ExpiresAt:               token.ExpiresAt.AsTime(),
		Encrypted:               token.Encrypted,
		InternalEncryptionLevel: int(token.InternalEncryptionLevel),
	}
}
