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
	"k8s.io/utils/strings/slices"
)

/*
	SendVerificationEmail - this method sends a verification text to the user with a code used to verify the phone.
*/
func (h *defaultHandler) SendVerificationText(ctx context.Context, req *go_hera.HeraRequest) (resp *go_hera.HeraResponse, err error) {
	var (
		userRepository   user_repository.UserRepository
		configRepository config_repository.ConfigRepository
		user             *go_hera.User
		randomCode       string
		verificationCode []byte
		config           *go_hera.Config
		errGroup         = &errgroup.Group{}
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
	// async action 2 - send verification email
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
	if user.GetPhone() == "" {
		return nil, errors.New("user do not have a phone - set the phone for the user")
	}
	if slices.Contains(user.VerifiedPhoneNumbers, user.PhoneHash) {
		return nil, errors.New("email is already verified")
	}
	// async action 1 - update phone verification code
	errGroup.Go(func() (err error) {
		if err = h.text.SendVerificationText(config.GetName(), user.GetPhone(), string(verificationCode)); err != nil {
			return err
		}
		return nil
	})
	// async action 2  update verification email sent at
	errGroup.Go(func() (err error) {
		if err = userRepository.UpdatePhoneVerificationCode(ctx, &go_hera.User{
			PhoneVerificationCode: &go_hera.Hash{Body: string(verificationCode)},
			Id:                    user.Id,
		}); err != nil {
			return err
		}
		return
	})
	if err = errGroup.Wait(); err != nil {
		return nil, err
	}
	return
}
