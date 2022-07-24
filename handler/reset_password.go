package handler

import (
	"context"
	"errors"
	"github.com/nuntiodev/hera-sdks/go_hera"
	"github.com/nuntiodev/hera/hash"
	"github.com/nuntiodev/hera/helpers"
	"github.com/nuntiodev/hera/repository/config_repository"
	"github.com/nuntiodev/hera/repository/user_repository"
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
		hasher           hash.Hash
		config           *go_hera.Config
		user             *go_hera.User
		hashErr          error
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
	hasher = hash.New(config)
	errGroup.Go(func() (err error) {
		// get requested user and check if the email is already verified
		userRepository, err = h.repository.UserRepositoryBuilder().SetNamespace(req.Namespace).SetHasher(hasher).WithPasswordValidation(config.ValidatePassword).Build(ctx)
		if err != nil {
			return err
		}
		user, err = userRepository.Get(ctx, req.User)
		if err != nil {
			return err
		}
		if user.ResetPasswordCode == nil || user.ResetPasswordCode.Body == "" {
			return errors.New("reset password code has not been set")
		}
		if req.User.ResetPasswordCode == nil || req.User.ResetPasswordCode.Body == "" {
			return errors.New("missing reset password code")
		}
		if time.Now().Sub(user.ResetPasswordCodeSentAt.AsTime()).Minutes() > h.maxVerificationAge.Minutes() {
			return errors.New("verification code has expired, send a new")
		}
		return nil
	})
	if err = errGroup.Wait(); err != nil {
		return nil, err
	}
	// provide exponential backoff
	time.Sleep(helpers.GetExponentialBackoff(float64(user.VerifyEmailAttempts), helpers.BackoffFactorTwo))
	hashErr = hasher.Compare(strings.TrimSpace(req.User.ResetPasswordCode.Body), user.ResetPasswordCode)
	if err = userRepository.UpdatePassword(ctx, user, req.User); err != nil {
		return nil, err
	}
	if hashErr != nil {
		return nil, err
	}
	return &go_hera.HeraResponse{}, nil
}
