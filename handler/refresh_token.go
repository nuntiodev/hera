package handler

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/nuntiodev/hera/repository/token_repository"
	"golang.org/x/sync/errgroup"
	"time"

	"github.com/nuntiodev/hera-proto/go_hera"
	"github.com/nuntiodev/hera/token"
	ts "google.golang.org/protobuf/types/known/timestamppb"
)

/*
	RefreshToken - this method provides a new access / refresh token pair given a valid refresh token.
*/
func (h *defaultHandler) RefreshToken(ctx context.Context, req *go_hera.HeraRequest) (resp *go_hera.HeraResponse, err error) {
	var (
		tokenRepository token_repository.TokenRepository
		refreshClaims   *go_hera.CustomClaims
		refreshToken    string
		accessClaims    *go_hera.CustomClaims
		accessToken     string
		errGroup        = &errgroup.Group{}
	)
	// async action 1 - validate that the refresh token is signed by Nuntio.
	errGroup.Go(func() (err error) {
		refreshClaims, err = h.token.ValidateToken(publicKey, req.Token.RefreshToken)
		if err != nil {
			return err
		}
		if refreshClaims.Type != token.TokenTypeRefresh {
			return errors.New("invalid refresh token")
		}
		return nil
	})
	// async action 2 - check if the refresh token is blocked.
	errGroup.Go(func() (err error) {
		// check if token is blocked in db
		tokenRepository, err = h.repository.TokenRepositoryBuilder().SetNamespace(req.Namespace).Build(ctx)
		if err != nil {
			return err
		}
		isBlocked, err := tokenRepository.IsBlocked(ctx, &go_hera.Token{
			Id:     refreshClaims.Id,
			UserId: refreshClaims.UserId,
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
	// create new refresh token and block the old one
	// if refresh token is about to expire (in less than 10 hours),
	refreshToken = req.Token.RefreshToken
	if _, err := h.BlockToken(ctx, &go_hera.HeraRequest{
		Token: &go_hera.Token{
			RefreshToken: refreshToken,
		},
	}); err != nil {
		return nil, err
	}
	refreshToken, refreshClaims, err = h.token.GenerateToken(privateKey, uuid.NewString(), refreshClaims.UserId, "", token.TokenTypeRefresh, refreshTokenExpiry)
	if err != nil {
		return nil, err
	}
	// create refresh token in database
	if err := tokenRepository.Create(ctx, &go_hera.Token{
		Id:           refreshClaims.Id,
		UserId:       refreshClaims.UserId,
		Type:         go_hera.TokenType_TOKEN_TYPE_REFRESH,
		LoggedInFrom: req.Token.LoggedInFrom,
		DeviceInfo:   req.Token.DeviceInfo,
		ExpiresAt:    ts.New(time.Unix(refreshClaims.ExpiresAt, 0)),
	}); err != nil {
		return nil, err
	}
	// generate new access token from refresh token
	accessToken, accessClaims, err = h.token.GenerateToken(privateKey, uuid.NewString(), refreshClaims.UserId, refreshClaims.Id, token.TokenTypeAccess, accessTokenExpiry)
	if err != nil {
		return nil, err
	}
	// async action 3 - add new access token info to database
	errGroup.Go(func() (err error) {
		if err = tokenRepository.Create(ctx, &go_hera.Token{
			Id:           accessClaims.Id,
			UserId:       accessClaims.UserId,
			Type:         go_hera.TokenType_TOKEN_TYPE_ACCESS,
			LoggedInFrom: req.Token.LoggedInFrom,
			DeviceInfo:   req.Token.DeviceInfo,
			ExpiresAt:    ts.New(time.Unix(accessClaims.ExpiresAt, 0)),
		}); err != nil {
			return err
		}
		return
	})
	// async action 4 - set refresh token used at to now
	errGroup.Go(func() (err error) {
		if err = tokenRepository.UpdateUsedAt(ctx, &go_hera.Token{
			Id:           refreshClaims.Id,
			UserId:       refreshClaims.UserId,
			LoggedInFrom: req.Token.LoggedInFrom,
			DeviceInfo:   req.Token.DeviceInfo,
		}); err != nil {
			return err
		}
		return err
	})
	err = errGroup.Wait()
	if err != nil {
		return nil, err
	}
	return &go_hera.HeraResponse{
		Token: &go_hera.Token{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	}, nil
}
