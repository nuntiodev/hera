package handler

import (
	"context"
	"errors"
	"github.com/nuntiodev/hera-sdks/go_hera"
	"github.com/nuntiodev/hera/helpers"
	"github.com/nuntiodev/hera/models"
	"github.com/nuntiodev/hera/repository/user_repository"
	"golang.org/x/crypto/bcrypt"
	"k8s.io/utils/strings/slices"
	"strings"
	"time"
)

/*
	VerifyEmail - this method verifies the code sent to a users email and sets the email as verified if correct.
*/
func (h *defaultHandler) VerifyEmail(ctx context.Context, req *go_hera.HeraRequest) (resp *go_hera.HeraResponse, err error) {
	var (
		userRepository user_repository.UserRepository
		user           *models.User
		bcryptErr      error
	)
	// get requested user and check if the email is already verified
	userRepository, err = h.repository.UserRepositoryBuilder().SetNamespace(req.Namespace).Build(ctx)
	if err != nil {
		return nil, err
	}
	user, err = userRepository.Get(ctx, req.User)
	if err != nil {
		return nil, err
	}
	if slices.Contains(user.VerifiedEmails, user.EmailHash) {
		return nil, errors.New("email is already verified")
	}
	if user.EmailVerificationCode == "" {
		return nil, errors.New("verification email has not been sent")
	}
	if req.User.EmailVerificationCode == "" {
		return nil, errors.New("missing provided email verification code")
	}
	if time.Now().Sub(user.VerificationEmailSentAt).Minutes() > h.maxVerificationAge.Minutes() {
		return nil, errors.New("verification email has expired, send a new one or login again")
	}
	// provide exponential backoff
	time.Sleep(helpers.GetExponentialBackoff(float64(user.VerifyEmailAttempts), helpers.BackoffFactorTwo))
	bcryptErr = bcrypt.CompareHashAndPassword([]byte(user.EmailVerificationCode), []byte(strings.TrimSpace(req.User.EmailVerificationCode)))
	// verify email
	if err = userRepository.VerifyEmail(ctx, models.UserToProtoUser(user), bcryptErr == nil); err != nil {
		return nil, err
	}
	if bcryptErr != nil {
		return nil, err
	}
	return &go_hera.HeraResponse{}, nil
}
