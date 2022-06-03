package handler

import (
	"context"
	"fmt"
	"github.com/nuntiodev/block-proto/go_block"
	"github.com/nuntiodev/nuntio-user-block/repository/config_repository"
	"github.com/nuntiodev/nuntio-user-block/repository/email_repository"
	"github.com/nuntiodev/nuntio-user-block/repository/user_repository"
	"golang.org/x/sync/errgroup"
)

/*
	CreateNamespaceConfig - this method creates a config for a given namespace.
*/
func (h *defaultHandler) CreateNamespaceConfig(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error) {
	var (
		userRepo   user_repository.UserRepository
		configRepo config_repository.ConfigRepository
		emailRepo  email_repository.EmailRepository
		config     *go_block.Config
		errGroup   = &errgroup.Group{}
		err        error
	)
	// async action 1 - setup user repository and create test user
	errGroup.Go(func() error {
		userRepo, err = h.repository.Users().SetNamespace(req.Namespace).Build(ctx)
		if err != nil {
			return fmt.Errorf("could not build user repository with err: %v", err)
		}
		// create test user
		_, err = userRepo.Create(ctx, &go_block.User{
			FirstName: "Test",
			LastName:  "User",
			Email:     "test@user.io",
		})
		return err
	})
	// async action 2 - setup config repository and create default config
	errGroup.Go(func() error {
		// create initial config
		configRepo, err = h.repository.Config(ctx, req.Namespace)
		if err != nil {
			return fmt.Errorf("could not build config with err: %v", err)
		}
		config, err = configRepo.Create(ctx, req.Config)
		return err
	})
	if err = errGroup.Wait(); err != nil {
		return &go_block.UserResponse{}, err
	}
	// create default emails
	emailRepo, err = h.repository.Email(ctx, req.Namespace)
	if err != nil {
		return &go_block.UserResponse{}, err
	}
	// async action 3 - create default verification email
	errGroup.Go(func() error {
		_, err = emailRepo.Create(ctx, &go_block.Email{
			Id:             email_repository.VerificationEmail,
			Logo:           config.Logo,
			WelcomeMessage: "Hello",
			BodyMessage:    fmt.Sprintf("Thank you for signing up to %s. In order to get started, we ask of you to confirm your email by entering the following numbers in your %s app.", config.Name, config.Name),
			FooterMessage:  fmt.Sprintf("All the best from %s team", config.Name),
			Title:          "Verify your email",
			Subject:        "Verify your email account",
			TemplatePath:   emailVerificationTemplatePath,
		})
		return err
	})
	// async action 4 - create default reset password email
	errGroup.Go(func() error {
		_, err = emailRepo.Create(ctx, &go_block.Email{
			Id:             email_repository.ResetPasswordEmail,
			Logo:           config.Logo,
			WelcomeMessage: "Hello",
			BodyMessage:    fmt.Sprintf("Thank you for signing up to %s. In order to reset your password, enter the following numbers in your %s app together with your new password.", config.Name, config.Name),
			FooterMessage:  fmt.Sprintf("All the best from %s team", config.Name),
			Title:          "Verify your email",
			Subject:        "Verify your email account",
			TemplatePath:   emailVerificationTemplatePath,
		})
		return err
	})
	if err = errGroup.Wait(); err != nil {
		return &go_block.UserResponse{}, nil
	}
	return &go_block.UserResponse{
		Config: config,
	}, nil
}
