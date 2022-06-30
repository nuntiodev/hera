package handler

import (
	"context"
	"errors"
	"github.com/nuntiodev/hera/repository/user_repository"

	"github.com/nuntiodev/hera-sdks/go_hera"
	"golang.org/x/crypto/bcrypt"
)

/*
	ValidateCredentials - this method is used to validate the credentials of a user without issuing a token.
*/
func (h *defaultHandler) ValidateCredentials(ctx context.Context, req *go_hera.HeraRequest) (resp *go_hera.HeraResponse, err error) {
	var (
		userRepository user_repository.UserRepository
		user           *go_hera.User
	)
	userRepository, err = h.repository.UserRepositoryBuilder().SetNamespace(req.Namespace).Build(ctx)
	if err != nil {
		return nil, err
	}
	user, err = userRepository.Get(ctx, req.User)
	if err != nil {
		return nil, err
	}
	if user.Password == "" {
		return nil, errors.New("please update the user with a non-empty password")
	}
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.User.Password)); err != nil {
		return nil, err
	}
	return &go_hera.HeraResponse{
		User: user,
	}, nil
}
