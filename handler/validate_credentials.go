package handler

import (
	"context"
	"errors"
	"github.com/nuntiodev/hera-sdks/go_hera"
	"github.com/nuntiodev/hera/hash"
	"github.com/nuntiodev/hera/repository/user_repository"
)

/*
	ValidateCredentials - this method is used to validate the credentials of a user without issuing a token.
*/
func (h *defaultHandler) ValidateCredentials(ctx context.Context, req *go_hera.HeraRequest) (resp *go_hera.HeraResponse, err error) {
	var (
		userRepository user_repository.UserRepository
		user           *go_hera.User
	)
	// async action 2 - get default config
	userRepository, err = h.repository.UserRepositoryBuilder().SetNamespace(req.Namespace).Build(ctx)
	if err != nil {
		return nil, err
	}
	user, err = userRepository.Get(ctx, req.User)
	if err != nil {
		return nil, err
	}
	if user.Password == nil || user.Password.Body == "" {
		return nil, errors.New("please update the user with a non-empty password")
	}
	if err = hash.New(nil).Compare(req.User.Password.Body, user.Password); err != nil {
		return nil, err
	}
	return &go_hera.HeraResponse{
		User: user,
	}, nil
}
