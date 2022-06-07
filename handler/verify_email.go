package handler

import (
	"context"
	"errors"
	"github.com/nuntiodev/block-proto/go_block"
	"github.com/nuntiodev/nuntio-user-block/helpers"
	"github.com/nuntiodev/nuntio-user-block/models"
	"github.com/nuntiodev/nuntio-user-block/repository/user_repository"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

/*
	VerifyEmail - this method verifies the code sent to a users email and sets the email as verified if correct.
*/
func (h *defaultHandler) VerifyEmail(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error) {
	var (
		userRepo user_repository.UserRepository
		user     *models.User
		err      error
	)
	// get requested user and check if the email is already verified
	userRepo, err = h.repository.UserRepositoryBuilder().SetNamespace(req.Namespace).SetEncryptionKey(req.EncryptionKey).Build(ctx)
	if err != nil {
		return &go_block.UserResponse{}, err
	}
	user, err = userRepo.Get(ctx, req.User)
	if err != nil {
		return &go_block.UserResponse{}, err
	}
	if user.EmailIsVerified {
		return &go_block.UserResponse{}, errors.New("email is already verified")
	}
	if user.EmailVerificationCode == "" {
		return &go_block.UserResponse{}, errors.New("verification email has not been sent")
	}
	if req.EmailVerificationCode == "" {
		return &go_block.UserResponse{}, errors.New("missing provided email verification code")
	}
	if time.Now().Sub(user.VerificationEmailSentAt).Minutes() > h.maxEmailVerificationAge.Minutes() {
		return &go_block.UserResponse{}, errors.New("verification email has expired, send a new one or login again")
	}
	// provide exponential backoff
	time.Sleep(helpers.GetExponentialBackoff(float64(user.VerifyEmailAttempts), helpers.BackoffFactorTwo))
	bcryptErr := bcrypt.CompareHashAndPassword([]byte(user.EmailVerificationCode), []byte(strings.TrimSpace(req.EmailVerificationCode)))
	if _, err := userRepo.UpdateEmailVerified(ctx, models.UserToProtoUser(user), &go_block.User{
		EmailIsVerified: bcryptErr == nil,
		EmailHash:       user.EmailHash,
	}); err != nil {
		return &go_block.UserResponse{}, err
	}
	return &go_block.UserResponse{}, bcryptErr
}
