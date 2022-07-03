package handler

import (
	"context"
	"errors"
	"github.com/nuntiodev/hera-sdks/go_hera"
	"github.com/nuntiodev/hera/repository/config_repository"
	"github.com/nuntiodev/hera/repository/user_repository"
)

/*
	CreateUser - this method creates a user in the database with a valid config.
*/
func (h *defaultHandler) CreateUser(ctx context.Context, req *go_hera.HeraRequest) (resp *go_hera.HeraResponse, err error) {
	var (
		configRepository config_repository.ConfigRepository
		userRepository   user_repository.UserRepository
		config           *go_hera.Config
		user             *go_hera.User
	)
	configRepository, err = h.repository.ConfigRepositoryBuilder().SetNamespace(req.Namespace).Build(ctx)
	if err != nil {
		return nil, err
	}
	config, err = configRepository.Get(ctx)
	if err != nil {
		return nil, err
	}
	// validate that the action is possible with project config
	// we cannot send an email if the email provider is not enabled
	if h.emailEnabled == false && config.VerifyEmail {
		return nil, errors.New("email provider is not enabled and verification email cannot be sent. If you want to enable email verification, override the EmailSender interface")
	} else if config.VerifyEmail && (req.User.Email == nil || *req.User.Email == "") {
		return nil, errors.New("require email verification is enabled and user email is empty. If you require email verification, specify an email for the user")
	}
	// we cannot send a text message if the text provider is enabled
	if !h.textEnabled && config.VerifyPhone {
		return nil, errors.New("text provider is not enabled and verification text cannot be sent. If you want to enable text verification, override the TextSender interface")
	} else if config.VerifyPhone && (req.User.GetPhone() == "") {
		return nil, errors.New("require phone number verification is enabled and user phone is empty. If you require phone number verification, specify a phone number for the user")
	}
	// build repository used to create user in the database
	h.repository.UserRepositoryBuilder()
	userRepository, err = h.repository.UserRepositoryBuilder().SetNamespace(req.Namespace).WithPasswordValidation(config.ValidatePassword).Build(ctx)
	if err != nil {
		return nil, err
	}
	user, err = userRepository.Create(ctx, req.User)
	if err != nil {
		return nil, err
	}
	return &go_hera.HeraResponse{
		User: user,
	}, nil
}
