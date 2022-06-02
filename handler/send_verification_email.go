package handler

import (
	"context"
	"encoding/hex"
	"errors"
	"github.com/nuntiodev/block-proto/go_block"
	"github.com/nuntiodev/nuntio-user-block/email"
	"github.com/nuntiodev/nuntio-user-block/repository/email_repository"
	"github.com/nuntiodev/nuntio-user-block/repository/user_repository"
	"github.com/nuntiodev/x/cryptox"
	"golang.org/x/sync/errgroup"
	"strings"
)

/*
	SendVerificationEmail - this method sends a verification email to the user with a code used to verify the email.
*/
func (h *defaultHandler) SendVerificationEmail(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error) {
	var (
		userRepo          user_repository.UserRepository
		emailRepo         email_repository.EmailRepository
		user              *go_block.User
		nameOfUser        string
		verificationCode  []byte
		verificationEmail *go_block.Email
		errGroup          = &errgroup.Group{}
		err               error
	)
	if !h.emailEnabled {
		return nil, errors.New("email provider is not enabled")
	}
	// async action 1 - get user and check if his email is verified
	errGroup.Go(func() error {
		userRepo, err = h.repository.Users().SetNamespace(req.Namespace).SetEncryptionKey(req.EncryptionKey).Build(ctx)
		if err != nil {
			return err
		}
		user, err = userRepo.Get(ctx, req.User, true)
		if err != nil {
			return err
		}
		if user.Email == "" {
			return errors.New("user do not have an email - set the email for the user")
		}
		if user.EmailIsVerified {
			return errors.New("email is already verified")
		}
		nameOfUser = user.Email
		if user.FirstName != "" {
			nameOfUser = strings.TrimSpace(user.FirstName + " " + user.LastName)
		}
		return err
	})
	// async action 2 - setup email repository and generate verification code
	errGroup.Go(func() error {
		emailRepo, err = h.repository.Email(ctx, req.Namespace)
		if err != nil {
			return err
		}
		randomCode, err := h.crypto.GenerateSymmetricKey(6, cryptox.Numeric)
		if err != nil {
			return err
		}
		verificationCode, err = hex.DecodeString(randomCode)
		if err != nil {
			return err
		}
		verificationEmail, err = emailRepo.Get(ctx, &go_block.Email{
			Id: email_repository.VerificationEmail,
		})
		return err
	})
	if err = errGroup.Wait(); err != nil {
		return &go_block.UserResponse{}, err
	}
	// async action 3 - send verification email
	errGroup.Go(func() error {
		return h.email.SendVerificationEmail(user.Email, verificationEmail.Subject, verificationEmail.TemplatePath, &email.VerificationData{
			Code: string(verificationCode),
			TemplateData: email.TemplateData{
				LogoUrl:        verificationEmail.Logo,
				WelcomeMessage: verificationEmail.WelcomeMessage,
				NameOfUser:     nameOfUser,
				BodyMessage:    verificationEmail.BodyMessage,
				FooterMessage:  verificationEmail.FooterMessage,
			},
		})
	})
	// async action 4  update verification email sent at
	errGroup.Go(func() error {
		_, err = userRepo.UpdateVerificationEmailSent(ctx, &go_block.User{
			EmailVerificationCode: string(verificationCode),
			Id:                    user.Id,
		})
		return err
	})
	err = errGroup.Wait()
	return &go_block.UserResponse{
		User: user,
	}, err
}
