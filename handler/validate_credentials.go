package handler

import (
	"context"
	"errors"
	"github.com/softcorp-io/block-proto/go_block"
	"golang.org/x/crypto/bcrypt"
)

func (h *defaultHandler) ValidateCredentials(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error) {
	resp, err := h.Get(ctx, req)
	if err != nil {
		return &go_block.UserResponse{}, err
	}
	if resp.User.Password == "" {
		return &go_block.UserResponse{}, errors.New("please update the user with a non-empty password")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(resp.User.Password), []byte(req.User.Password)); err != nil {
		return &go_block.UserResponse{}, err
	}
	return &go_block.UserResponse{
		User: resp.User,
	}, nil
}
