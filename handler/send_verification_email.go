package handler

import (
	"context"
	"errors"
	"github.com/nuntiodev/hera-proto/go_hera"
	"github.com/nuntiodev/hera/models"
	"github.com/nuntiodev/hera/repository/user_repository"
	"golang.org/x/sync/errgroup"
	"k8s.io/utils/strings/slices"
	"strings"
)

/*
	SendVerificationEmail - this method sends a verification email to the user with a code used to verify the email.
*/
func (h *defaultHandler) SendVerificationEmail(ctx context.Context, req *go_hera.HeraRequest) (resp *go_hera.HeraResponse, err error) {
	var (
		userRepository   user_repository.UserRepository
		user             *models.User
		nameOfUser       string
		verificationCode []byte
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
		if slices.Contains(user.VerifiedEmails, user.EmailHash) {
			return errors.New("email is already verified")
		}
		nameOfUser = user.Email.Body
		if user.FirstName.Body != "" {
			nameOfUser = strings.TrimSpace(user.FirstName.Body)
			if user.LastName.Body != "" {
				nameOfUser += " " + user.LastName.Body
			}
		}
		return
	})
	if err = errGroup.Wait(); err != nil {
		return nil, err
	}
	// async action 2 - send verification email
	errGroup.Go(func() (err error) {
		if err = h.email.SendVerificationEmail(user.Email.Body, string(verificationCode)); err != nil {
			return err
		}
		return
	})
	// async action 3  update verification email sent at
	errGroup.Go(func() (err error) {
		if err = userRepository.UpdateEmailVerificationCode(ctx, &go_hera.User{
			EmailVerificationCode: string(verificationCode),
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
		User: models.UserToProtoUser(user),
	}, nil
}
