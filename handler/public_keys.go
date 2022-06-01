package handler

import (
	"context"

	"github.com/nuntiodev/block-proto/go_block"
)

/*
	PublicKeys - this method returns all internal public keys.
*/
func (h *defaultHandler) PublicKeys(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error) {
	return &go_block.UserResponse{
		PublicKeys: map[string]string{
			"public-jwt-key": publicKeyString,
		},
	}, nil
}
