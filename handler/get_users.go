package handler

import (
	"context"
	"github.com/nuntiodev/hera/models"
	"github.com/nuntiodev/hera/repository/user_repository"

	"github.com/nuntiodev/hera-sdks/go_hera"
)

/*
	GetUsers - this method fetches a collection of users from the database.
*/
func (h *defaultHandler) GetUsers(ctx context.Context, req *go_hera.HeraRequest) (resp *go_hera.HeraResponse, err error) {
	var (
		userRepository user_repository.UserRepository
		users          []*models.User
	)
	userRepository, err = h.repository.UserRepositoryBuilder().SetNamespace(req.Namespace).Build(ctx)
	if err != nil {
		return nil, err
	}
	users, err = userRepository.GetMany(ctx, req.Users)
	if err != nil {
		return nil, err
	}
	return &go_hera.HeraResponse{
		Users: models.UsersToProto(users),
	}, nil
}
