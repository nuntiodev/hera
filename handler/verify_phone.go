package handler

import (
	"context"
	"errors"
	"github.com/nuntiodev/hera-proto/go_hera"
	"github.com/nuntiodev/hera/helpers"
	"github.com/nuntiodev/hera/models"
	"github.com/nuntiodev/hera/repository/user_repository"
	"golang.org/x/crypto/bcrypt"
	"k8s.io/utils/strings/slices"
	"strings"
	"time"
)

/*
	VerifyPhone - this method verifies the code sent to a users phone and sets the phone as verified if correct.
*/
func (h *defaultHandler) VerifyPhone(ctx context.Context, req *go_hera.HeraRequest) (resp *go_hera.HeraResponse, err error) {
	var (
		userRepository user_repository.UserRepository
		user           *models.User
		bcryptErr      error
	)
	// get requested user and check if the phone is already verified
	userRepository, err = h.repository.UserRepositoryBuilder().SetNamespace(req.Namespace).Build(ctx)
	if err != nil {
		return nil, err
	}
	user, err = userRepository.Get(ctx, req.User)
	if err != nil {
		return nil, err
	}
	if slices.Contains(user.VerifiedPhoneNumbers, user.PhoneHash) {
		return nil, errors.New("phone is already verified")
	}
	if user.PhoneVerificationCode == "" {
		return nil, errors.New("verification text has not been sent")
	}
	if req.User.PhoneVerificationCode == "" {
		return nil, errors.New("missing provided text verification code")
	}
	if time.Now().Sub(user.VerificationTextSentAt).Minutes() > h.maxVerificationAge.Minutes() {
		return nil, errors.New("verification text has expired, send a new one or login again")
	}
	// provide exponential backoff
	time.Sleep(helpers.GetExponentialBackoff(float64(user.VerifyPhoneAttempts), helpers.BackoffFactorTwo))
	bcryptErr = bcrypt.CompareHashAndPassword([]byte(user.PhoneVerificationCode), []byte(strings.TrimSpace(req.User.PhoneVerificationCode)))
	// verify phone
	if err = userRepository.VerifyPhone(ctx, models.UserToProtoUser(user), bcryptErr == nil); err != nil {
		return nil, err
	}
	if bcryptErr != nil {
		return nil, err
	}
	return &go_hera.HeraResponse{}, nil
}
