package handler

import (
	"context"
	"encoding/hex"
	"errors"
	"github.com/nuntiodev/hera-sdks/go_hera"
	"github.com/nuntiodev/hera/hash"
	"github.com/nuntiodev/hera/repository/config_repository"
	"github.com/nuntiodev/hera/repository/user_repository"
	"github.com/nuntiodev/x/cryptox"
	"golang.org/x/sync/errgroup"
)

/*
	SendResetPasswordEmail - this method send an email verification code to the user.
*/
func (h *defaultHandler) SendResetPasswordEmail(ctx context.Context, req *go_hera.HeraRequest) (resp *go_hera.HeraResponse, err error) {
	var (
		userRepository   user_repository.UserRepository
		configRepository config_repository.ConfigRepository
		user             *go_hera.User
		randomCode       string
		verificationCode []byte
		config           *go_hera.Config
		errGroup         = errgroup.Group{}
	)
	if h.emailEnabled == false {
		return nil, errors.New("email provider is not enabled")
	}
	randomCode, err = cryptox.GenerateSymmetricKey(6, cryptox.Numeric)
	if err != nil {
		return nil, err
	}
	verificationCode, err = hex.DecodeString(randomCode)
	if err != nil {
		return nil, err
	}
	configRepository, err = h.repository.ConfigRepositoryBuilder().SetNamespace(req.Namespace).Build(ctx)
	if err != nil {
		return nil, err
	}
	config, err = configRepository.Get(ctx)
	if err != nil {
		return nil, err
	}
	userRepository, err = h.repository.UserRepositoryBuilder().SetHasher(hash.New(config)).SetNamespace(req.Namespace).Build(ctx)
	if err != nil {
		return nil, err
	}
	user, err = userRepository.Get(ctx, req.User)
	if err != nil {
		return nil, err
	}
	if user.GetEmail() == "" {
		return nil, errors.New("user do not have an email - set the email for the user")
	}
	// async action 1 - send reset password email
	errGroup.Go(func() (err error) {
		if err = h.email.SendResetPasswordEmail(config.GetName(), user.GetEmail(), string(verificationCode)); err != nil {
			return err
		}
		return nil
	})
	// async action 2 - update reset password code
	errGroup.Go(func() (err error) {
		if err = userRepository.UpdateResetPasswordCode(ctx, &go_hera.User{
			EmailVerificationCode: &go_hera.Hash{Body: string(verificationCode)},
			Id:                    user.Id,
		}); err != nil {
			return err
		}
		return nil
	})
	if err = errGroup.Wait(); err != nil {
		return nil, err
	}
	return
}
