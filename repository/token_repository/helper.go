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
		Id:                      token.Id,
		UserId:                  token.UserId,
		Blocked:                 token.Blocked,
		DeviceInfo:              token.Device,
		LoggedInFrom:            token.LoggedInFrom,
		Type:                    token.Type,
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
	return &Token{
		Id:                      token.Id,
		UserId:                  token.UserId,
		Blocked:                 token.Blocked,
		Device:                  token.DeviceInfo,
		LoggedInFrom:            token.LoggedInFrom,
		Type:                    token.Type,
		BlockedAt:               token.BlockedAt.AsTime(),
		CreatedAt:               token.CreatedAt.AsTime(),
		UsedAt:                  token.UsedAt.AsTime(),
		ExpiresAt:               token.ExpiresAt.AsTime(),
		Encrypted:               token.Encrypted,
		InternalEncryptionLevel: int(token.InternalEncryptionLevel),
	}
}
