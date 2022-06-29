package handler

import (
	"context"
	"github.com/nuntiodev/hera/repository/user_repository"

	"github.com/nuntiodev/hera-sdks/go_hera"
)

/*
	DeleteUsers - this method deletes a batch of users.
*/
func (h *defaultHandler) DeleteUsers(ctx context.Context, req *go_hera.HeraRequest) (resp *go_hera.HeraResponse, err error) {
	var (
		userRepository user_repository.UserRepository
	)
	userRepository, err = h.repository.UserRepositoryBuilder().SetNamespace(req.Namespace).Build(ctx)
	if err != nil {
		return nil, err
	}
	return nil, userRepository.DeleteMany(ctx, req.Users)
}
