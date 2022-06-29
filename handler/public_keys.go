package handler

import (
	"context"

	"github.com/nuntiodev/hera-sdks/go_hera"
)

/*
	PublicKeys - this method returns all internal public keys.
*/
func (h *defaultHandler) PublicKeys(ctx context.Context, req *go_hera.HeraRequest) (resp *go_hera.HeraResponse, err error) {
	return &go_hera.HeraResponse{
		PublicKeys: map[string]string{
			"hera-public-key": publicKeyString,
		},
	}, nil
}
