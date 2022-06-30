package handler

import (
	"context"
	"encoding/hex"
	"errors"
	"github.com/nuntiodev/hera-sdks/go_hera"
	"github.com/nuntiodev/hera/repository/config_repository"
	"github.com/nuntiodev/hera/repository/user_repository"
	"github.com/nuntiodev/x/cryptox"
	"golang.org/x/sync/errgroup"
	"strings"
)

/*
	SendResetPasswordText - this method send a text verification code to the user.
*/
func (h *defaultHandler) SendResetPasswordText(ctx context.Context, req *go_hera.HeraRequest) (resp *go_hera.HeraResponse, err error) {
	var (
		userRepository   user_repository.UserRepository
		configRepository config_repository.ConfigRepository
		user             *go_hera.User
		nameOfUser       string
		randomCode       string
		verificationCode []byte
		config           *go_hera.Config
		errGroup         = &errgroup.Group{}
	)
	if h.textEnabled == false {
		return nil, errors.New("text provider is not enabled")
	}
	// async action 1 - get user and check if his phone is verified
	errGroup.Go(func() (err error) {
		userRepository, err = h.repository.UserRepositoryBuilder().SetNamespace(req.Namespace).Build(ctx)
		if err != nil {
			return err
		}
		user, err = userRepository.Get(ctx, req.User)
		if err != nil {
			return err
		}
		if user.GetPhone() == "" {
			return errors.New("user do not have a phone number - set the phone number for the user")
		}
		nameOfUser = user.GetPhone()
		if user.GetFirstName() != "" {
			nameOfUser = strings.TrimSpace(user.GetFirstName())
			if user.GetLastName() != "" {
				nameOfUser += " " + user.GetLastName()
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
	// async action 2 - send verification text
	errGroup.Go(func() (err error) {
		configRepository, err = h.repository.ConfigRepositoryBuilder().SetNamespace(req.Namespace).Build(ctx)
		if err != nil {
			return err
		}
		config, err = configRepository.Get(ctx)
		if err != nil {
			return err
		}
		if err = h.text.SendResetPasswordText(config.GetName(), user.GetPhone(), string(verificationCode)); err != nil {
			return err
		}
		return err
	})
	// async action 3  update verification phone sent at
	errGroup.Go(func() (err error) {
		if err = userRepository.UpdateResetPasswordCode(ctx, &go_hera.User{
			PhoneVerificationCode: string(verificationCode),
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
