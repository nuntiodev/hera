package handler

import (
	"context"

	"github.com/io-nuntio/block-proto/go_block"
)

func (h *defaultHandler) PublicKeys(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error) {
	return &go_block.UserResponse{
		PublicKeys: map[string]string{
			"public-jwt-key": publicKeyString,
		},
	}, nil
}
