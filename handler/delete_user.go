package handler

import (
	"context"
	"github.com/nuntiodev/hera/repository/user_repository"

	"github.com/nuntiodev/hera-proto/go_hera"
)

/*
	DeleteUser - this method deletes a user from a given repository.
*/
func (h *defaultHandler) DeleteUser(ctx context.Context, req *go_hera.HeraRequest) (resp *go_hera.HeraResponse, err error) {
	var (
		userRepository user_repository.UserRepository
	)
	userRepository, err = h.repository.UserRepositoryBuilder().SetNamespace(req.Namespace).Build(ctx)
	if err != nil {
		return nil, err
	}
	if err = userRepository.Delete(ctx, req.User); err != nil {
		return nil, err
	}
	return &go_hera.HeraResponse{}, nil
}
