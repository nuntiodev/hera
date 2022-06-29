package handler

import (
	"context"
	"errors"
	"github.com/nuntiodev/hera-sdks/go_hera"
	"github.com/nuntiodev/hera/helpers"
	"github.com/nuntiodev/hera/models"
	"github.com/nuntiodev/hera/repository/config_repository"
	"github.com/nuntiodev/hera/repository/user_repository"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/sync/errgroup"
	"strings"
	"time"
)

/*
	ResetPassword - this method validates that the provided verification code matches the hashed code stored in the database
	and updates the users password.
*/
func (h *defaultHandler) ResetPassword(ctx context.Context, req *go_hera.HeraRequest) (resp *go_hera.HeraResponse, err error) {
	var (
		configRepository config_repository.ConfigRepository
		userRepository   user_repository.UserRepository
		config           *models.Config
		user             *models.User
		bcryptErr        error
		errGroup         errgroup.Group
	)
	configRepository, err = h.repository.ConfigRepositoryBuilder().SetNamespace(req.Namespace).Build(ctx)
	if err != nil {
		return nil, err
	}
	config, err = configRepository.Get(ctx)
	if err != nil {
		return nil, err
	}
	errGroup.Go(func() (err error) {
		// get requested user and check if the email is already verified
		userRepository, err = h.repository.UserRepositoryBuilder().SetNamespace(req.Namespace).WithPasswordValidation(config.ValidatePassword).Build(ctx)
		if err != nil {
			return err
		}
		user, err = userRepository.Get(ctx, req.User)
		if err != nil {
			return err
		}
		if user.ResetPasswordCode == "" {
			return errors.New("reset password code has not been set")
		}
		if req.User.ResetPasswordCode == "" {
			return errors.New("missing reset password code")
		}
		if time.Now().Sub(user.ResetPasswordCodeSentAt).Minutes() > h.maxVerificationAge.Minutes() {
			return errors.New("verification code has expired, send a new")
		}
		return nil
	})
	if err = errGroup.Wait(); err != nil {
		return nil, err
	}
	// provide exponential backoff
	time.Sleep(helpers.GetExponentialBackoff(float64(user.VerifyEmailAttempts), helpers.BackoffFactorTwo))
	bcryptErr = bcrypt.CompareHashAndPassword([]byte(user.ResetPasswordCode), []byte(strings.TrimSpace(req.User.ResetPasswordCode)))
	// reset password
	if err = userRepository.UpdatePassword(ctx, models.UserToProtoUser(user), req.User); err != nil {
		return nil, err
	}
	if bcryptErr != nil {
		return nil, err
	}
	return &go_hera.HeraResponse{}, nil
}
