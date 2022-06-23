package handler

import (
	"context"
	"github.com/nuntiodev/hera/repository/user_repository"

	"github.com/nuntiodev/hera-proto/go_hera"
)

/*
	UpdateUserMetadata - this method updates a users metadata which is stored as JSON and used to store additional
	information about a user.
*/
func (h *defaultHandler) UpdateUserMetadata(ctx context.Context, req *go_hera.HeraRequest) (resp *go_hera.HeraResponse, err error) {
	var (
		userRepository user_repository.UserRepository
	)
	userRepository, err = h.repository.UserRepositoryBuilder().SetNamespace(req.Namespace).Build(ctx)
	if err != nil {
		return nil, err
	}
	if err = userRepository.UpdateMetadata(ctx, req.User, req.UserUpdate); err != nil {
		return nil, err
	}
	return nil, nil
}
