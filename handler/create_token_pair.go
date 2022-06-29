package handler

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/nuntiodev/hera/models"
	"github.com/nuntiodev/hera/repository/user_repository"
	"golang.org/x/sync/errgroup"
	"time"

	"github.com/nuntiodev/hera-sdks/go_hera"
	"github.com/nuntiodev/hera/token"
	ts "google.golang.org/protobuf/types/known/timestamppb"
)

/*
	CreateTokenPair - this method create an access and refresh token for a given user.
*/
func (h *defaultHandler) CreateTokenPair(ctx context.Context, req *go_hera.HeraRequest) (resp *go_hera.HeraResponse, err error) {
	var (
		userRepository user_repository.UserRepository
		user           *models.User
		refreshToken   string
		refreshClaims  *go_hera.CustomClaims
		accessToken    string
		accessClaims   *go_hera.CustomClaims
		errGroup       = &errgroup.Group{}
	)
	userRepository, err = h.repository.UserRepositoryBuilder().SetNamespace(req.Namespace).Build(ctx)
	if err != nil {
		return nil, err
	}
	user, err = userRepository.Get(ctx, req.User)
	if err != nil {
		return nil, err
	}
	//  generate and save refresh and access tokenRepository
	tokenRepository, err := h.repository.TokenRepositoryBuilder().SetNamespace(req.Namespace).Build(ctx)
	if err != nil {
		return nil, fmt.Errorf("could setup token repository with err: %v", err)
	}
	// build data for tokenRepository
	refreshTokenId := uuid.NewString()
	loggedInFrom := ""
	deviceInfo := ""
	if req.Token != nil {
		loggedInFrom = req.Token.LoggedInFrom
		deviceInfo = req.Token.DeviceInfo
	}
	// async action 1 - generate refresh token.
	errGroup.Go(func() (err error) {
		refreshToken, refreshClaims, err = h.token.GenerateToken(privateKey, refreshTokenId, user.Id, "", token.RefreshToken, refreshTokenExpiry)
		if err != nil {
			return fmt.Errorf("could generate refresh token with err: %v", err)
		}
		// create refresh token info in database
		if err = tokenRepository.Create(ctx, &go_hera.Token{
			Id:           refreshClaims.Id,
			UserId:       refreshClaims.UserId,
			ExpiresAt:    ts.New(time.Unix(refreshClaims.ExpiresAt, 0)),
			LoggedInFrom: loggedInFrom,
			DeviceInfo:   deviceInfo,
			Type:         go_hera.TokenType_TOKEN_TYPE_REFRESH,
		}); err != nil {
			return err
		}
		return nil
	})
	// async action 2 - generate access token.
	errGroup.Go(func() (err error) {
		accessToken, accessClaims, err = h.token.GenerateToken(privateKey, uuid.NewString(), user.Id, refreshTokenId, token.AccessToken, accessTokenExpiry)
		if err != nil {
			return fmt.Errorf("could generate access token with err: %v", err)
		}
		// create access token info in database
		if err = tokenRepository.Create(ctx, &go_hera.Token{
			Id:           accessClaims.Id,
			UserId:       accessClaims.UserId,
			ExpiresAt:    ts.New(time.Unix(accessClaims.ExpiresAt, 0)),
			LoggedInFrom: loggedInFrom,
			DeviceInfo:   deviceInfo,
			Type:         go_hera.TokenType_TOKEN_TYPE_ACCESS,
		}); err != nil {
			return err
		}
		return nil
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
		User: models.UserToProtoUser(user),
	}, nil
}
