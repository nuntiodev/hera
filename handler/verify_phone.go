package handler

import (
	"context"
	"errors"
	"github.com/nuntiodev/hera-sdks/go_hera"
	"github.com/nuntiodev/hera/hash"
	"github.com/nuntiodev/hera/helpers"
	"github.com/nuntiodev/hera/repository/user_repository"
	"k8s.io/utils/strings/slices"
	"time"
)

/*
	VerifyPhone - this method verifies the code sent to a users phone and sets the phone as verified if correct.
*/
func (h *defaultHandler) VerifyPhone(ctx context.Context, req *go_hera.HeraRequest) (resp *go_hera.HeraResponse, err error) {
	var (
		userRepository user_repository.UserRepository
		user           *go_hera.User
		hashErr        error
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
	if user.PhoneVerificationCode == nil || user.PhoneVerificationCode.Body == "" {
		return nil, errors.New("verification text has not been sent")
	}
	if req.User.PhoneVerificationCode == nil || req.User.PhoneVerificationCode.Body == "" {
		return nil, errors.New("missing provided text verification code")
	}
	if time.Now().Sub(user.VerificationTextSentAt.AsTime()).Minutes() > h.maxVerificationAge.Minutes() {
		return nil, errors.New("verification text has expired, send a new one or login again")
	}
	// provide exponential backoff
	time.Sleep(helpers.GetExponentialBackoff(float64(user.VerifyPhoneAttempts), helpers.BackoffFactorTwo))
	hashErr = hash.New(nil).Compare(req.User.PhoneVerificationCode.Body, user.PhoneVerificationCode)
	// verify phone
	if err = userRepository.VerifyPhone(ctx, user, hashErr == nil); err != nil {
		return nil, err
	}
	if hashErr != nil {
		return nil, err
	}
	return &go_hera.HeraResponse{}, nil
}
