package handler

import (
	"context"
	"github.com/nuntiodev/hera/models"
	"github.com/nuntiodev/hera/repository/user_repository"

	"github.com/nuntiodev/hera-proto/go_hera"
)

/*
	SearchForUser - this method fetches a user from the database by querying.
*/
func (h *defaultHandler) SearchForUser(ctx context.Context, req *go_hera.HeraRequest) (resp *go_hera.HeraResponse, err error) {
	var (
		userRepository user_repository.UserRepository
		user           *models.User
	)
	userRepository, err = h.repository.UserRepositoryBuilder().SetNamespace(req.Namespace).Build(ctx)
	if err != nil {
		return nil, err
	}
	user, err = userRepository.Search(ctx, req.Query.Search)
	if err != nil {
		return nil, err
	}
	return &go_hera.HeraResponse{
		User: models.UserToProtoUser(user),
	}, nil
}
