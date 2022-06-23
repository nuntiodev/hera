package handler

import (
	"context"
	"encoding/hex"
	"errors"
	"github.com/nuntiodev/hera-proto/go_hera"
	"github.com/nuntiodev/hera/models"
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
		user             *models.User
		nameOfUser       string
		randomCode       string
		verificationCode []byte
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
		if user.Phone.Body == "" {
			return errors.New("user do not have a phone number - set the phone number for the user")
		}
		nameOfUser = user.Phone.Body
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
	// async action 2 - send verification text
	errGroup.Go(func() (err error) {
		if err = h.text.SendResetPasswordText(user.Phone.Body, string(verificationCode)); err != nil {
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
