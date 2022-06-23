package handler

import (
	"context"
	"errors"
	"github.com/nuntiodev/hera-proto/go_hera"
	"github.com/nuntiodev/hera/helpers"
	"github.com/nuntiodev/hera/models"
	"github.com/nuntiodev/hera/repository/user_repository"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

/*
	ResetPassword - this method validates that the provided verification code matches the hashed code stored in the database
	and updates the users password.
*/
func (h *defaultHandler) ResetPassword(ctx context.Context, req *go_hera.HeraRequest) (resp *go_hera.HeraResponse, err error) {
	var (
		userRepository user_repository.UserRepository
		user           *models.User
		bcryptErr      error
	)
	// get requested user and check if the email is already verified
	userRepository, err = h.repository.UserRepositoryBuilder().SetNamespace(req.Namespace).Build(ctx)
	if err != nil {
		return nil, err
	}
	user, err = userRepository.Get(ctx, req.User)
	if err != nil {
		return nil, err
	}
	if user.ResetPasswordCode == "" {
		return nil, errors.New("reset password code has not been set")
	}
	if req.User.ResetPasswordCode == "" {
		return nil, errors.New("missing reset password code")
	}
	if time.Now().Sub(user.ResetPasswordEmailSentAt).Minutes() > h.maxVerificationAge.Minutes() {
		return nil, errors.New("verification code has expired, send a new")
	}
	// provide exponential backoff
	time.Sleep(helpers.GetExponentialBackoff(float64(user.VerifyEmailAttempts), helpers.BackoffFactorTwo))
	bcryptErr = bcrypt.CompareHashAndPassword([]byte(user.ResetPasswordCode), []byte(strings.TrimSpace(req.User.ResetPasswordCode)))
	// verify email
	//todo: look through this
	if err = userRepository.UpdateResetPasswordCode(ctx, models.UserToProtoUser(req.User)); err != nil {
		return nil, err
	}
	if bcryptErr != nil {
		return nil, err
	}
	return nil, nil
}
