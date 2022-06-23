package handler

import (
	"context"
	"github.com/nuntiodev/hera-proto/go_hera"
)

/*
	ResetPassword - this method validates that the provided verification code matches the hashed code stored in the database
	and updates the users password.
	todo: enable both email and text reset password.
*/
func (h *defaultHandler) ResetPassword(ctx context.Context, req *go_hera.HeraRequest) (resp *go_hera.HeraResponse, err error) {
	return nil, nil
}
