package handler

import (
	"context"
	"errors"
	"github.com/nuntiodev/hera-sdks/go_hera"
	"github.com/nuntiodev/hera/repository/config_repository"
	"github.com/nuntiodev/hera/repository/user_repository"
	"golang.org/x/sync/errgroup"
	"k8s.io/utils/strings/slices"
	"strings"
)

/*
	SendVerificationText - this method sends a verification text to the user with a code used to verify the phone.
*/
func (h *defaultHandler) SendVerificationText(ctx context.Context, req *go_hera.HeraRequest) (resp *go_hera.HeraResponse, err error) {
	var (
		userRepository   user_repository.UserRepository
		configRepository config_repository.ConfigRepository
		user             *go_hera.User
		nameOfUser       string
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
		if slices.Contains(user.VerifiedPhoneNumbers, user.PhoneHash) {
			return errors.New("phone is already verified")
		}
		nameOfUser = user.GetPhone()
		if user.GetFirstName() != "" {
			nameOfUser = strings.TrimSpace(user.GetFirstName())
			if user.GetLastName() != "" {
				nameOfUser += " " + user.GetLastName()
			}
		}
		return
	})
	if err = errGroup.Wait(); err != nil {
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
		if err = h.text.SendVerificationText(config.GetName(), user.GetPhone(), string(verificationCode)); err != nil {
			return err
		}
		return
	})
	// async action 3  update verification text sent at
	errGroup.Go(func() (err error) {
		if err = userRepository.UpdatePhoneVerificationCode(ctx, &go_hera.User{
			PhoneVerificationCode: string(verificationCode),
			Id:                    user.Id,
		}); err != nil {
			return err
		}
		return
	})
	if err = errGroup.Wait(); err != nil {
		return nil, err
	}
	return &go_hera.HeraResponse{
		User: user,
	}, nil
}
