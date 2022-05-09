package handler

import (
	"context"
	"errors"
	"github.com/nuntiodev/block-proto/go_block"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

func (h *defaultHandler) VerifyEmail(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error) {
	// get requested user and check if the email is already verified
	userResp, err := h.Get(ctx, req)
	if err != nil {
		return &go_block.UserResponse{}, err
	}
	get := userResp.User
	if get.EmailIsVerified {
		return &go_block.UserResponse{}, errors.New("email is already verified")
	}
	if get.EmailVerificationCode == "" {
		return &go_block.UserResponse{}, errors.New("verification email has not been sent")
	}
	if req.EmailVerificationCode == "" {
		return &go_block.UserResponse{}, errors.New("missing provided email verification code")
	}
	if time.Now().UTC().Sub(get.VerificationEmailSentAt.AsTime()).Minutes() > maxEmailVerificationAge.Minutes() {
		return &go_block.UserResponse{}, errors.New("verification email has expired, send a new one or login again")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(get.EmailVerificationCode), []byte(strings.TrimSpace(req.EmailVerificationCode))); err != nil {
		return &go_block.UserResponse{}, err
	}
	users, err := h.repository.Users().SetNamespace(req.Namespace).Build(ctx)
	if err != nil {
		return &go_block.UserResponse{}, err
	}
	if _, err := users.UpdateEmailVerified(ctx, get, &go_block.User{
		EmailIsVerified: true,
	}); err != nil {
		return &go_block.UserResponse{}, err
	}
	return &go_block.UserResponse{}, nil
}
