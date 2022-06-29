package handler

import (
	"context"
	"github.com/nuntiodev/hera/repository/user_repository"

	"github.com/nuntiodev/hera-sdks/go_hera"
)

/*
	UpdateUserContact - this method updates a users contact information - this includes a users email, phone, and/or username.
*/
func (h *defaultHandler) UpdateUserContact(ctx context.Context, req *go_hera.HeraRequest) (resp *go_hera.HeraResponse, err error) {
	var (
		userRepository user_repository.UserRepository
	)
	userRepository, err = h.repository.UserRepositoryBuilder().SetNamespace(req.Namespace).Build(ctx)
	if err != nil {
		return nil, err
	}
	if err = userRepository.UpdateContact(ctx, req.User, req.UserUpdate); err != nil {
		return nil, err
	}
	return &go_hera.HeraResponse{}, nil
}
