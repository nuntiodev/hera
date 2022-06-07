package handler

import (
	"context"
	"encoding/hex"
	"errors"
	"github.com/nuntiodev/block-proto/go_block"
	"github.com/nuntiodev/nuntio-user-block/email"
	"github.com/nuntiodev/nuntio-user-block/models"
	"github.com/nuntiodev/nuntio-user-block/repository/email_repository"
	"github.com/nuntiodev/nuntio-user-block/repository/user_repository"
	"github.com/nuntiodev/x/cryptox"
	"golang.org/x/sync/errgroup"
	"strings"
)

/*
	SendResetPasswordEmail - this method send an email verification code to the user.
	todo: enable both email and text reset password.
*/
func (h *defaultHandler) SendResetPasswordEmail(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error) {
	var (
		userRepo          user_repository.UserRepository
		emailRepo         email_repository.EmailRepository
		user              *models.User
		nameOfUser        string
		verificationCode  []byte
		verificationEmail *models.Email
		errGroup          = &errgroup.Group{}
		err               error
	)
	if !h.emailEnabled {
		return nil, errors.New("email provider is not enabled")
	}
	// async action 1 - get user and check if his email is verified
	errGroup.Go(func() error {
		userRepo, err = h.repository.UserRepositoryBuilder().SetNamespace(req.Namespace).SetEncryptionKey(req.EncryptionKey).Build(ctx)
		if err != nil {
			return err
		}
		user, err = userRepo.Get(ctx, req.User)
		if err != nil {
			return err
		}
		if user.Email.Body == "" {
			return errors.New("user do not have an email - set the email for the user")
		}
		if user.EmailIsVerified {
			return errors.New("email is already verified")
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
	// async action 2 - setup email repository and generate verification code
	errGroup.Go(func() error {
		emailRepo, err = h.repository.Email(ctx, req.Namespace, req.EncryptionKey)
		if err != nil {
			return err
		}
		randomCode, err := cryptox.GenerateSymmetricKey(6, cryptox.Numeric)
		if err != nil {
			return err
		}
		verificationCode, err = hex.DecodeString(randomCode)
		if err != nil {
			return err
		}
		verificationEmail, err = emailRepo.Get(ctx, &go_block.Email{
			Id: email_repository.ResetPasswordEmail,
		})
		return err
	})
	if err = errGroup.Wait(); err != nil {
		return &go_block.UserResponse{}, err
	}
	// async action 3 - send verification email
	errGroup.Go(func() error {
		return h.email.SendVerificationEmail(user.Email.Body, verificationEmail.Subject.Body, verificationEmail.TemplatePath.Body, &email.VerificationData{
			Code: string(verificationCode),
			TemplateData: email.TemplateData{
				LogoUrl:        verificationEmail.Logo.Body,
				WelcomeMessage: verificationEmail.WelcomeMessage.Body,
				NameOfUser:     nameOfUser,
				BodyMessage:    verificationEmail.BodyMessage.Body,
				FooterMessage:  verificationEmail.FooterMessage.Body,
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
	return &go_block.UserResponse{}, err
}
