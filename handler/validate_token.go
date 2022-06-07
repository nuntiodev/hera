package handler

import (
	"context"
	"errors"
	"github.com/nuntiodev/nuntio-user-block/models"
	"github.com/nuntiodev/nuntio-user-block/repository/token_repository"
	"github.com/nuntiodev/nuntio-user-block/repository/user_repository"
	"golang.org/x/sync/errgroup"

	"github.com/nuntiodev/block-proto/go_block"
	"github.com/nuntiodev/nuntio-user-block/token"
)

/*
	ValidateToken - this method validates a users token and stores information about it in the database.
*/
func (h *defaultHandler) ValidateToken(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error) {
	var (
		tokenRepo token_repository.TokenRepository
		userRepo  user_repository.UserRepository
		claims    *go_block.CustomClaims
		user      *models.User
		errGroup  = &errgroup.Group{}
		err       error
	)
	// async action 1 - validate token is issued by Nuntio
	errGroup.Go(func() error {
		claims, err = h.token.ValidateToken(publicKey, req.TokenPointer)
		return err
	})
	// async action 2 - build token repository
	errGroup.Go(func() error {
		tokenRepo, err = h.repository.Tokens(ctx, req.Namespace, req.EncryptionKey)
		return err
	})
	if err = errGroup.Wait(); err != nil {
		return &go_block.UserResponse{}, err
	}
	// validate if token is blocked in db
	// async action 3 - validate if access token is blocked
	errGroup.Go(func() error {
		// for access tokens we also validate if refresh token is blocked
		if claims.Type == token.TokenTypeAccess {
			isBlocked, err := tokenRepo.IsBlocked(ctx, &go_block.Token{
				Id:     claims.RefreshTokenId,
				UserId: claims.UserId,
			})
			if err != nil {
				return err
			}
			if isBlocked {
				return errors.New("token is blocked")
			}
		}
		return nil
	})
	errGroup.Go(func() error {
		// we always validate if the pointer token is blocked
		isBlocked, err := tokenRepo.IsBlocked(ctx, &go_block.Token{
			Id:     claims.Id,
			UserId: claims.UserId,
		})
		if err != nil {
			return err
		}
		if isBlocked {
			return errors.New("token is blocked")
		}
		return nil
	})
	if err = errGroup.Wait(); err != nil {
		return &go_block.UserResponse{}, err
	}
	// async action 4 - build data for token and store in database.
	errGroup.Go(func() error {
		loggedInFrom := &go_block.Location{}
		deviceInfo := ""
		if req.Token != nil {
			loggedInFrom = req.Token.LoggedInFrom
			deviceInfo = req.Token.DeviceInfo
		}
		_, err := tokenRepo.UpdateUsedAt(ctx, &go_block.Token{
			Id:           claims.Id,
			LoggedInFrom: loggedInFrom,
			DeviceInfo:   deviceInfo,
		})
		return err
	})
	// async action 5 - get user from token user id
	userRepo, err = h.repository.UserRepositoryBuilder().SetNamespace(req.Namespace).SetEncryptionKey(req.EncryptionKey).Build(ctx)
	if err != nil {
		return &go_block.UserResponse{}, err
	}
	user, err = userRepo.Get(ctx, &go_block.User{Id: claims.UserId})
	return &go_block.UserResponse{
		User: models.UserToProtoUser(user),
	}, err
}
