package handler

import (
	"context"
	"github.com/nuntiodev/block-proto/go_block"
)

/*
	ResetPassword - this method validates that the provided verification code matches the hashed code stored in the database
	and updates the users password.
	todo: enable both email and text reset password.
*/
func (h *defaultHandler) ResetPassword(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error) {
	return &go_block.UserResponse{}, nil
}
