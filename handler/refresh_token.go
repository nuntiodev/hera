package handler

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/nuntiodev/nuntio-user-block/repository/token_repository"
	"golang.org/x/sync/errgroup"
	"time"

	"github.com/nuntiodev/block-proto/go_block"
	"github.com/nuntiodev/nuntio-user-block/token"
	ts "google.golang.org/protobuf/types/known/timestamppb"
)

/*
	RefreshToken - this method provides a new access / refresh token pair given a valid refresh token.
*/
func (h *defaultHandler) RefreshToken(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error) {
	var (
		tokenRepo     token_repository.TokenRepository
		refreshClaims *go_block.CustomClaims
		refreshToken  string
		accessClaims  *go_block.CustomClaims
		accessToken   string
		errGroup      = &errgroup.Group{}
		err           error
	)
	// async action 1 - validate that the refresh token is signed by Nuntio.
	errGroup.Go(func() error {
		refreshClaims, err = h.token.ValidateToken(publicKey, req.Token.RefreshToken)
		if err != nil {
			return err
		}
		if refreshClaims.Type != token.TokenTypeRefresh {
			return errors.New("invalid refresh token")
		}
		return nil
	})
	// async action 2 - check if token is blocked.
	errGroup.Go(func() error {
		// check if token is blocked in db
		tokenRepo, err = h.repository.Tokens(ctx, req.Namespace)
		if err != nil {
			return err
		}
		isBlocked, err := tokenRepo.IsBlocked(ctx, &go_block.Token{
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
		return &go_block.UserResponse{}, err
	}
	// build metadata for token
	loggedInFrom := &go_block.Location{}
	deviceInfo := ""
	if req.Token != nil {
		loggedInFrom = req.Token.LoggedInFrom
		deviceInfo = req.Token.DeviceInfo
	}
	// create new refresh token if the token is expired.
	// if refresh token is about to expire (in less than 10 hours), create a new one and block the old one
	refreshToken = req.Token.RefreshToken
	if time.Unix(refreshClaims.ExpiresAt, 0).Sub(time.Now()) < time.Hour*10 {
		if _, err := h.BlockToken(ctx, &go_block.UserRequest{
			Token: &go_block.Token{
				RefreshToken: refreshToken,
			},
		}); err != nil {
			return &go_block.UserResponse{}, err
		}
		refreshToken, refreshClaims, err = h.token.GenerateToken(privateKey, uuid.NewString(), refreshClaims.UserId, "", token.TokenTypeRefresh, refreshTokenExpiry)
		if err != nil {
			return &go_block.UserResponse{}, err
		}
		// create refresh token in database
		if _, err := tokenRepo.Create(ctx, &go_block.Token{
			Id:           refreshClaims.Id,
			UserId:       refreshClaims.UserId,
			Type:         go_block.TokenType_TOKEN_TYPE_REFRESH,
			LoggedInFrom: loggedInFrom,
			DeviceInfo:   deviceInfo,
			ExpiresAt:    ts.New(time.Unix(refreshClaims.ExpiresAt, 0)),
		}); err != nil {
			return &go_block.UserResponse{}, err
		}
	}
	// generate new access token from refresh token
	accessToken, accessClaims, err = h.token.GenerateToken(privateKey, uuid.NewString(), refreshClaims.UserId, refreshClaims.Id, token.TokenTypeAccess, accessTokenExpiry)
	if err != nil {
		return &go_block.UserResponse{}, err
	}
	// async action 3 - add new access token info to database
	errGroup.Go(func() error {
		_, err := tokenRepo.Create(ctx, &go_block.Token{
			Id:           accessClaims.Id,
			UserId:       accessClaims.UserId,
			Type:         go_block.TokenType_TOKEN_TYPE_ACCESS,
			LoggedInFrom: loggedInFrom,
			DeviceInfo:   deviceInfo,
			ExpiresAt:    ts.New(time.Unix(accessClaims.ExpiresAt, 0)),
		})
		return err
	})
	// async action 4 - set refresh token used at to now
	errGroup.Go(func() error {
		_, err = tokenRepo.UpdateUsedAt(ctx, &go_block.Token{
			Id:           refreshClaims.Id,
			UserId:       refreshClaims.UserId,
			LoggedInFrom: loggedInFrom,
			DeviceInfo:   deviceInfo,
		})
		return err
	})
	err = errGroup.Wait()
	return &go_block.UserResponse{
		Token: &go_block.Token{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	}, err
}
