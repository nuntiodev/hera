package handler

import (
	"context"
	"github.com/nuntiodev/hera/models"
	"github.com/nuntiodev/hera/repository/user_repository"

	"github.com/nuntiodev/hera-proto/go_hera"
)

/*
	UpdateUserProfile - this method updates a users profile - this includes a users first name, last name, image and/or birthdate.
*/
func (h *defaultHandler) UpdateUserProfile(ctx context.Context, req *go_hera.HeraRequest) (resp *go_hera.HeraResponse, err error) {
	var (
		userRepository user_repository.UserRepository
		user           *models.User
	)
	userRepository, err = h.repository.UserRepositoryBuilder().SetNamespace(req.Namespace).Build(ctx)
	if err != nil {
		return nil, err
	}
	if err = userRepository.UpdateProfile(ctx, req.User, req.UserUpdate); err != nil {
		return nil, err
	}
	return &go_hera.HeraResponse{
		User: models.UserToProtoUser(user),
	}, nil
}
