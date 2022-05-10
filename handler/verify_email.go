package handler

import (
	"context"
	"errors"
	"github.com/nuntiodev/block-proto/go_block"
	"github.com/nuntiodev/nuntio-user-block/helpers"
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
	if time.Now().Sub(get.VerificationEmailSentAt.AsTime()).Minutes() > maxEmailVerificationAge.Minutes() {
		return &go_block.UserResponse{}, errors.New("verification email has expired, send a new one or login again")
	}
	// provide exponential backoff
	time.Sleep(helpers.GetExponentialBackoff(float64(get.VerifyEmailAttempts)))
	bcryptErr := bcrypt.CompareHashAndPassword([]byte(get.EmailVerificationCode), []byte(strings.TrimSpace(req.EmailVerificationCode)))
	users, err := h.repository.Users().SetNamespace(req.Namespace).Build(ctx)
	if err != nil {
		return &go_block.UserResponse{}, err
	}
	if _, err := users.UpdateEmailVerified(ctx, get, &go_block.User{
		EmailIsVerified: bcryptErr == nil,
		EmailHash:       get.EmailHash,
	}); err != nil {
		return &go_block.UserResponse{}, err
	}
	if bcryptErr != nil {
		return &go_block.UserResponse{}, bcryptErr
	}
	return &go_block.UserResponse{}, nil
}
