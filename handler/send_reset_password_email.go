package handler

import (
	"context"
	"encoding/hex"
	"errors"
	"github.com/nuntiodev/hera-proto/go_hera"
	"github.com/nuntiodev/hera/models"
	"github.com/nuntiodev/hera/repository/config_repository"
	"github.com/nuntiodev/hera/repository/user_repository"
	"github.com/nuntiodev/x/cryptox"
	"golang.org/x/sync/errgroup"
	"strings"
)

/*
	SendResetPasswordEmail - this method send an email verification code to the user.
*/
func (h *defaultHandler) SendResetPasswordEmail(ctx context.Context, req *go_hera.HeraRequest) (resp *go_hera.HeraResponse, err error) {
	var (
		userRepository   user_repository.UserRepository
		configRepository config_repository.ConfigRepository
		user             *models.User
		nameOfUser       string
		randomCode       string
		verificationCode []byte
		config           *models.Config
		errGroup         = &errgroup.Group{}
	)
	if h.emailEnabled == false {
		return nil, errors.New("email provider is not enabled")
	}
	// async action 1 - get user and check if his email is verified
	errGroup.Go(func() (err error) {
		userRepository, err = h.repository.UserRepositoryBuilder().SetNamespace(req.Namespace).Build(ctx)
		if err != nil {
			return err
		}
		user, err = userRepository.Get(ctx, req.User)
		if err != nil {
			return err
		}
		if user.Email.Body == "" {
			return errors.New("user do not have an email - set the email for the user")
		}
		nameOfUser = user.Email.Body
		if user.FirstName.Body != "" {
			nameOfUser = strings.TrimSpace(user.FirstName.Body)
			if user.LastName.Body != "" {
				nameOfUser += " " + user.LastName.Body
			}
		}
		return err
	})
	if err = errGroup.Wait(); err != nil {
		return nil, err
	}
	randomCode, err = cryptox.GenerateSymmetricKey(6, cryptox.Numeric)
	if err != nil {
		return nil, err
	}
	verificationCode, err = hex.DecodeString(randomCode)
	if err != nil {
		return nil, err
	}
	// async action 2 - send verification email
	errGroup.Go(func() (err error) {
		configRepository, err = h.repository.ConfigRepositoryBuilder().SetNamespace(req.Namespace).Build(ctx)
		if err != nil {
			return err
		}
		config, err = configRepository.Get(ctx)
		if err != nil {
			return err
		}
		if err = h.email.SendResetPasswordEmail(config.Name.Body, user.Email.Body, string(verificationCode)); err != nil {
			return err
		}
		return err
	})
	// async action 3  update verification email sent at
	errGroup.Go(func() (err error) {
		if err = userRepository.UpdateResetPasswordCode(ctx, &go_hera.User{
			EmailVerificationCode: string(verificationCode),
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
