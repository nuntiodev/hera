package handler

import (
	"context"
	"errors"
	"github.com/nuntiodev/hera/repository/token_repository"
	"golang.org/x/sync/errgroup"

	"github.com/nuntiodev/hera-sdks/go_hera"
	"github.com/nuntiodev/hera/token"
)

/*
	ValidateToken - this method validates a users token and stores information about it in the database.
*/
func (h *defaultHandler) ValidateToken(ctx context.Context, req *go_hera.HeraRequest) (resp *go_hera.HeraResponse, err error) {
	var (
		tokenRepository token_repository.TokenRepository
		claims          *go_hera.CustomClaims
		errGroup        = &errgroup.Group{}
	)
	// async action 1 - validate token is issued by Nuntio
	errGroup.Go(func() (err error) {
		claims, err = h.token.ValidateToken(publicKey, req.TokenPointer)
		if err != nil {
			return err
		}
		return nil
	})
	// async action 2 - build token repository
	errGroup.Go(func() (error error) {
		tokenRepository, err = h.repository.TokenRepositoryBuilder().SetNamespace(req.Namespace).Build(ctx)
		if err != nil {
			return err
		}
		return nil
	})
	if err = errGroup.Wait(); err != nil {
		return nil, err
	}
	// validate if token is blocked in db
	// async action 3 - validate if access token is blocked
	errGroup.Go(func() (err error) {
		// for access tokens we also validate if refresh token is blocked
		if claims.Type == token.AccessToken {
			var isBlocked bool
			isBlocked, err = tokenRepository.IsBlocked(ctx, &go_hera.Token{
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
	errGroup.Go(func() (error error) {
		var isBlocked bool
		// we always validate if the pointer token is blocked
		isBlocked, err = tokenRepository.IsBlocked(ctx, &go_hera.Token{
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
		return nil, err
	}
	// async action 4 - build data for token and store in database.
	errGroup.Go(func() (error error) {
		loggedInFrom := ""
		deviceInfo := ""
		if req.Token != nil {
			loggedInFrom = req.Token.LoggedInFrom
			deviceInfo = req.Token.DeviceInfo
		}
		if err = tokenRepository.UpdateUsedAt(ctx, &go_hera.Token{
			Id:           claims.Id,
			LoggedInFrom: loggedInFrom,
			DeviceInfo:   deviceInfo,
		}); err != nil {
			return err
		}
		return nil
	})
	return &go_hera.HeraResponse{}, nil
}
